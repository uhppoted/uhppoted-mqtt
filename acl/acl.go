package acl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"github.com/uhppoted/uhppoted-mqtt/log"
)

const (
	StatusInternalServerError = uhppoted.StatusInternalServerError
	StatusBadRequest          = uhppoted.StatusBadRequest
)

type Verification int

const (
	None Verification = iota
	NotEmpty
	RSA
)

func (v Verification) String() string {
	return [...]string{"none", "not-empty", "RSA"}[v]
}

type ACL struct {
	UHPPOTE     uhppote.IUHPPOTE
	Devices     []uhppote.Device
	RSA         *auth.RSA
	Credentials *credentials.Credentials
	Region      string
	Verify      map[Verification]bool
}

type Permission struct {
	Door      string     `json:"door"`
	StartDate types.Date `json:"start-date"`
	EndDate   types.Date `json:"end-date"`
	Profile   int        `json:"profile,omitempty"`
}

type Permissions struct {
	CardNumber  uint32       `json:"card-number"`
	Permissions []Permission `json:"permissions"`
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

func (a *ACL) fetch(tag, uri string, mimetype string) (*api.ACL, error) {
	infof(tag, "Fetching ACL from %v", uri)

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

	infof(tag, "Fetched ACL from %v (%d bytes)", uri, len(b))

	var extract func(io.Reader) (map[string][]byte, string, error)

	switch {
	case strings.HasSuffix(uri, ".zip"): // FIXME: remove - superseded by mime-type
		extract = unzip

	case mimetype == "application/zip":
		extract = unzip

	case mimetype == "text/tab-separated-values":
		extract = unpackTSV

	default:
		extract = untar
	}

	files, uname, err := extract(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	tsv, ok := files["ACL"]
	if !ok {
		return nil, fmt.Errorf("ACL file missing from tar.gz")
	}

	// ... verify signature if RSA or default verification
	signed := false

	if signature, ok := files["signature"]; !ok {
		warnf(tag, "'signature' file missing from tar.gz")
	} else {
		infof(tag, "Extracted ACL from %v: %v bytes, signature: %v bytes", uri, len(tsv), len(signature))

		if err := a.verify(uname, tsv, signature); err != nil {
			warnf(tag, "%v", err)
		} else {
			infof(tag, "ACL file '%v': verified signature", uri)
			signed = true
		}
	}

	acl, _, err := api.ParseTSV(bytes.NewReader(tsv), a.Devices, true)
	if err != nil {
		return nil, err
	}

	count := 0
	for _, v := range acl {
		count += len(v)
	}

	// ... verify
	switch {
	case a.Verify[None]:
		// 'k, good to go

	case a.Verify[NotEmpty] && !a.Verify[RSA]:
		if count == 0 {
			return nil, fmt.Errorf("ACL file %q: no records", uri)
		}

	case a.Verify[NotEmpty] && a.Verify[RSA]:
		if count == 0 && !signed {
			return nil, fmt.Errorf("ACL file %q': no records", uri)
		} else if count == 0 && signed {
			warnf(tag, "ACL file '%v':signed but contains no records", uri)
		}

	case a.Verify[RSA]:
		if !signed {
			return nil, fmt.Errorf("ACL file %q: invalid signature", uri)
		}

	default:
		if !signed {
			return nil, fmt.Errorf("ACL file %q: invalid signature", uri)
		}
	}

	infof(tag, "ACL file %q: verified", uri)

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

	infof(tag, "tar'd ACL (%v bytes) and signature (%v bytes): %v bytes", len(files["uhppoted.acl"]), len(files["signature"]), b.Len())

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

	infof(tag, "INFO  Stored ACL to %v", uri)

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

func infof(tag, format string, args ...any) {
	log.Infof(tag, format, args...)
}

func warnf(tag, format string, args ...any) {
	log.Warnf(tag, format, args...)
}
