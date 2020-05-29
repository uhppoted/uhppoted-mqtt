package acl

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type ACL struct {
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

type Error struct {
	Code    int    `json:"error-code"`
	Message string `json:"message"`
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

func fetchHTTP(url string) ([]byte, error) {
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

func fetchS3(url string, credentials *credentials.Credentials, region string) ([]byte, error) {
	match := regexp.MustCompile("^s3://(.*?)/(.*)").FindStringSubmatch(url)
	if len(match) != 3 {
		return nil, fmt.Errorf("Invalid S3 URI (%s)", url)
	}

	bucket := match[1]
	key := match[2]
	object := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	cfg := aws.NewConfig().
		WithCredentials(credentials).
		WithRegion(region)

	ss := session.Must(session.NewSession(cfg))

	buffer := make([]byte, 1024)
	b := aws.NewWriteAtBuffer(buffer)
	if _, err := s3manager.NewDownloader(ss).Download(b, &object); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func fetchFile(url string) ([]byte, error) {
	match := regexp.MustCompile("^file://(.*)").FindStringSubmatch(url)
	if len(match) != 2 {
		return nil, fmt.Errorf("Invalid file URI (%s)", url)
	}

	return ioutil.ReadFile(match[1])
}

func storeHTTP(uri string, r io.Reader) error {
	rq, err := http.NewRequest("PUT", "http://localhost:8080/upload", r)
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

func storeS3(uri string, credentials *credentials.Credentials, region string, r io.Reader) error {
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
		WithCredentials(credentials).
		WithRegion(region)

	ss := session.Must(session.NewSession(cfg))
	if _, err := s3manager.NewUploader(ss).Upload(&object); err != nil {
		return err
	}

	return nil
}

func storeFile(url string, r io.Reader) error {
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
