package acl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/auth"
)

const (
	StatusInternalServerError = uhppoted.StatusInternalServerError
	StatusBadRequest          = uhppoted.StatusBadRequest
)

type ACL struct {
	Devices     []*uhppote.Device
	RSA         *auth.RSA
	Credentials *credentials.Credentials
	Region      string
	Log         *log.Logger
	NoVerify    bool
}

type Permission struct {
	Door      string     `json:"door"`
	StartDate types.Date `json:"start-date"`
	EndDate   types.Date `json:"end-date"`
}

type Permissions struct {
	CardNumber  uint32       `json:"card-number"`
	Permissions []Permission `json:"permissions"`
}

func (a *ACL) info(tag, msg string) {
	a.Log.Printf("INFO  %-12s %s", tag, msg)
}

func (a *ACL) sign(acl []byte) ([]byte, error) {
	if a.RSA != nil {
		return a.RSA.Sign(acl)
	}

	return nil, nil
}

func (a *ACL) verify(uname string, acl, signature []byte) error {
	if a.RSA != nil {
		return a.RSA.Validate(uname, acl, signature)
	}

	return nil
}

func (a *ACL) fetch(tag, uri string) (*api.ACL, error) {
	a.info(tag, fmt.Sprintf("Fetching ACL from %v", uri))

	f := a.fetchHTTP
	if strings.HasPrefix(uri, "s3://") {
		f = a.fetchS3
	} else if strings.HasPrefix(uri, "file://") {
		f = a.fetchFile
	}

	b, err := f(uri)
	if err != nil {
		return nil, err
	}

	a.info(tag, fmt.Sprintf("Fetched ACL from %v (%d bytes)", uri, len(b)))

	x := untar
	if strings.HasSuffix(uri, ".zip") {
		x = unzip
	}

	files, uname, err := x(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	tsv, ok := files["ACL"]
	if !ok {
		return nil, fmt.Errorf("ACL file missing from tar.gz")
	}

	signature, ok := files["signature"]
	if !a.NoVerify && !ok {
		return nil, fmt.Errorf("'signature' file missing from tar.gz")
	}

	a.info(tag, fmt.Sprintf("Extracted ACL from %v: %v bytes, signature: %v bytes", uri, len(tsv), len(signature)))

	if !a.NoVerify {
		if err := a.verify(uname, tsv, signature); err != nil {
			return nil, err
		}
	}

	acl, _, err := api.ParseTSV(bytes.NewReader(tsv), a.Devices, true)
	if err != nil {
		return nil, err
	}

	return &acl, nil
}

func (a *ACL) fetchHTTP(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var b bytes.Buffer
	if _, err = io.Copy(&b, response.Body); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a *ACL) fetchS3(uri string) ([]byte, error) {
	match := regexp.MustCompile("^s3://(.*?)/(.*)").FindStringSubmatch(uri)
	if len(match) != 3 {
		return nil, fmt.Errorf("Invalid S3 URI (%s)", uri)
	}

	bucket := match[1]
	key := match[2]
	object := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	cfg := aws.NewConfig().
		WithCredentials(a.Credentials).
		WithRegion(a.Region)

	ss := session.Must(session.NewSession(cfg))

	buffer := make([]byte, 1024)
	b := aws.NewWriteAtBuffer(buffer)
	if _, err := s3manager.NewDownloader(ss).Download(b, &object); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a *ACL) fetchFile(url string) ([]byte, error) {
	match := regexp.MustCompile("^file://(.*)").FindStringSubmatch(url)
	if len(match) != 2 {
		return nil, fmt.Errorf("Invalid file URI (%s)", url)
	}

	return ioutil.ReadFile(match[1])
}

func (a *ACL) store(tag, uri, filename string, content []byte) error {
	files := map[string][]byte{
		filename: content,
	}

	if signature, err := a.sign(content); err != nil {
		return err
	} else if signature != nil {
		files["signature"] = signature
	}

	var b bytes.Buffer
	compress := targz
	if strings.HasSuffix(uri, ".zip") {
		compress = zipf
	}

	if err := compress(files, &b); err != nil {
		return err
	}

	a.info(tag, fmt.Sprintf("tar'd ACL (%v bytes) and signature (%v bytes): %v bytes", len(files["uhppoted.acl"]), len(files["signature"]), b.Len()))

	f := a.storeHTTP
	if strings.HasPrefix(uri, "s3://") {
		f = a.storeS3
	} else if strings.HasPrefix(uri, "file://") {
		f = a.storeFile
	} else {
	}

	if err := f(uri, bytes.NewReader(b.Bytes())); err != nil {
		return err
	}

	a.info(tag, fmt.Sprintf("INFO  Stored ACL to %v", uri))

	return nil
}

func (a *ACL) storeHTTP(url string, r io.Reader) error {
	rq, err := http.NewRequest("PUT", url, r)
	if err != nil {
		return err
	}

	rq.Header.Set("Content-Type", "binary/octet-stream")

	response, err := http.DefaultClient.Do(rq)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func (a *ACL) storeS3(uri string, r io.Reader) error {
	match := regexp.MustCompile("^s3://(.*?)/(.*)").FindStringSubmatch(uri)
	if len(match) != 3 {
		return fmt.Errorf("Invalid S3 URI (%s)", uri)
	}

	bucket := match[1]
	key := match[2]

	object := s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   r,
	}

	cfg := aws.NewConfig().
		WithCredentials(a.Credentials).
		WithRegion(a.Region)

	ss := session.Must(session.NewSession(cfg))
	if _, err := s3manager.NewUploader(ss).Upload(&object); err != nil {
		return err
	}

	return nil
}

func (a *ACL) storeFile(url string, r io.Reader) error {
	match := regexp.MustCompile("^file://(.*)").FindStringSubmatch(url)
	if len(match) != 2 {
		return fmt.Errorf("Invalid file URI (%s)", url)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(match[1], b, 0660)
}
