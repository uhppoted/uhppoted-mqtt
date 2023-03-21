package auth

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/uhppoted/uhppoted-mqtt/log"
)

type keyset struct {
	guard      sync.Mutex
	key        *rsa.PrivateKey
	clientKeys map[string]*rsa.PublicKey
	directory  string
}

type RSA struct {
	signingKeys    keyset
	encryptionKeys keyset
}

func NewRSA(keydir string) (*RSA, error) {
	var err error

	r := RSA{
		signingKeys: keyset{
			guard:      sync.Mutex{},
			key:        nil,
			clientKeys: map[string]*rsa.PublicKey{},
			directory:  filepath.Join(keydir, "signing"),
		},
		encryptionKeys: keyset{
			guard:      sync.Mutex{},
			key:        nil,
			clientKeys: map[string]*rsa.PublicKey{},
			directory:  filepath.Join(keydir, "encryption"),
		},
	}

	f := func(ks *keyset) error {
		keys, err := loadPublicKeys(ks.directory)
		if err != nil {
			return err
		}

		ks.guard.Lock()
		defer ks.guard.Unlock()

		ks.clientKeys = keys

		return nil
	}

	r.signingKeys.key, err = loadPrivateKey(filepath.Join(r.signingKeys.directory, "mqttd.key"))
	if err != nil {
		log.Warnf("RSA", "%v", err)
	}

	r.encryptionKeys.key, err = loadPrivateKey(filepath.Join(r.encryptionKeys.directory, "mqttd.key"))
	if err != nil {
		log.Warnf("RSA", "%v", err)
	}

	if err := f(&r.signingKeys); err != nil {
		log.Warnf("RSA", "%v", err)
	}

	if err := f(&r.encryptionKeys); err != nil {
		log.Warnf("RSA", "%v", err)
	}

	watch("signing keys", r.signingKeys.directory, func() error { return f(&r.signingKeys) })
	watch("encryption keys", r.encryptionKeys.directory, func() error { return f(&r.encryptionKeys) })

	return &r, nil
}

func (r *RSA) Validate(clientID string, request []byte, signature []byte) error {
	pubkey, ok := r.signingKeys.clientKeys[clientID]
	if !ok || pubkey == nil {
		return fmt.Errorf("no RSA public key (%v)", clientID)
	}

	hash := sha256.Sum256(request)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("%v  (%v)", err, clientID)
	}

	return nil
}

func (r *RSA) Sign(message []byte) ([]byte, error) {
	key := r.signingKeys.key
	if key != nil {
		rng := rand.Reader
		hashed := sha256.Sum256(message)

		return rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	}

	return []byte{}, nil
}

func (r *RSA) Encrypt(plaintext []byte, clientID string, label string) ([]byte, []byte, error) {
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, nil, err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	ciphertext := make([]byte, len(plaintext)+padding)
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...))

	rng := rand.Reader
	hash := sha256.Sum256([]byte(label))
	pubkey, ok := r.encryptionKeys.clientKeys[clientID]
	if !ok {
		return nil, nil, fmt.Errorf("no public key for %s", clientID)
	}

	key, err := rsa.EncryptOAEP(sha256.New(), rng, pubkey, secretKey, hash[:16])
	if err != nil {
		return nil, nil, err
	}

	return append(iv, ciphertext...), key, nil
}

func (r *RSA) Decrypt(ciphertext []byte, key []byte, label string) ([]byte, error) {
	rng := rand.Reader
	hash := sha256.Sum256([]byte(label))
	secretKey, err := rsa.DecryptOAEP(sha256.New(), rng, r.encryptionKeys.key, key, hash[:16])
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < 16 {
		return nil, fmt.Errorf("missing IV (%d bytes)", len(ciphertext))
	}

	if len(ciphertext[16:]) < aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext length (%d bytes)", len(ciphertext))
	}

	if len(ciphertext[16:])%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext not a multiple of AES block size (%d bytes)", len(ciphertext))
	}

	// REMOVED: using openssl AES on the command with the -salt option prepends the ciphertext with 'Salted__<salt>'
	// Ref. http://justsolve.archiveteam.org/wiki/OpenSSL_salted_format
	// offset := 0
	// if strings.HasPrefix(string(ciphertext), "Salted__") {
	// 	offset = 16
	// }
	//
	// plaintext := make([]byte, len(ciphertext[offset:]))
	// cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext[offset:])

	iv := ciphertext[:16]
	plaintext := make([]byte, len(ciphertext[16:]))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext[16:])

	N := len(plaintext)
	padding := int(plaintext[N-1])
	N -= padding

	if N < 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	return plaintext[:N], nil
}

func loadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	return pk, nil
}

func loadPublicKeys(dir string) (map[string]*rsa.PublicKey, error) {
	keys := map[string]*rsa.PublicKey{}
	filemode, err := os.Stat(dir)
	if err != nil {
		return keys, err
	}

	if !filemode.IsDir() {
		return keys, fmt.Errorf("%s is not a directory", dir)
	}

	list, err := os.ReadDir(dir)
	if err != nil {
		return keys, err
	}

	for _, f := range list {
		filename := f.Name()
		ext := filepath.Ext(filename)
		if ext == ".pub" {
			clientID := strings.TrimSuffix(filename, ext)

			bytes, err := os.ReadFile(filepath.Join(dir, filename))
			if err != nil {
				log.Warnf("RSA", "%v", err)
			}

			block, _ := pem.Decode(bytes)
			if block == nil || block.Type != "PUBLIC KEY" {
				log.Warnf("RSA", "%s is not a valid RSA public key", filename)
				continue
			}

			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				log.Warnf("RSA", "%s is not a valid RSA public key (%v)", filename, err)
				continue
			}

			pubkey, ok := key.(*rsa.PublicKey)
			if !ok {
				log.Warnf("RSA", "%s is not a valid RSA public key", filename)
				continue
			}

			keys[clientID] = pubkey
		}
	}

	return keys, nil
}

// NOTE: interim file watcher implementation pending fsnotify in Go 1.4
func watch(name string, directory string, reload func() error) {
	go func() {
		finfo, err := os.Stat(directory)
		if err != nil {
			log.Warnf("RSA", "Failed to get directory information for '%s': %v", directory, err)
			return
		}

		lastModified := finfo.ModTime()
		logged := false
		for {
			time.Sleep(2500 * time.Millisecond)
			finfo, err := os.Stat(directory)
			if err != nil {
				if !logged {
					log.Warnf("RSA", "Failed to get directory information for '%s': %v", directory, err)
					logged = true
				}

				continue
			}

			logged = false
			if finfo.ModTime() != lastModified {
				log.Infof("RSA", "Reloading information from %s\n", directory)

				err := reload()
				if err != nil {
					log.Errorf("RSA", "Failed to reload information from %s: %v", directory, err)
					continue
				}

				log.Warnf("RSA", "Updated %s from %s", name, directory)
				lastModified = finfo.ModTime()
			}
		}
	}()
}
