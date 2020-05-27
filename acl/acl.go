package acl

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/uhppoted/uhppote-core/types"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

func sign(acl []byte, keyfile string) ([]byte, error) {
	//	return auth.Sign(acl, keyfile)

	return nil, nil
}

func targz(files map[string][]byte, w io.Writer) error {
	var b bytes.Buffer

	tw := tar.NewWriter(&b)
	for filename, body := range files {
		header := &tar.Header{
			Name:  filename,
			Mode:  0660,
			Size:  int64(len(body)),
			Uname: "uhppoted",
			Gname: "uhppoted",
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := tw.Write([]byte(body)); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	gz := gzip.NewWriter(w)

	gz.Name = fmt.Sprintf("uhppoted-%s.tar.gz", time.Now().Format("2006-01-02T150405"))
	gz.ModTime = time.Now()
	gz.Comment = ""

	_, err := gz.Write(b.Bytes())
	if err != nil {
		return err
	}

	return gz.Close()
}

func zipf(files map[string][]byte, w io.Writer) error {
	zw := zip.NewWriter(w)
	for filename, body := range files {
		if f, err := zw.Create(filename); err != nil {
			return err
		} else if _, err = f.Write([]byte(body)); err != nil {
			return err
		}
	}

	return zw.Close()
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

func storeS3(uri, awsconfig, region string, r io.Reader) error {
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

	credentials, err := getAWSCredentials(awsconfig)
	if err != nil {
		return err
	}

	cfg := aws.NewConfig().
		WithCredentials(credentials).
		WithRegion(region)

	ss := session.Must(session.NewSession(cfg))

	_, err = s3manager.NewUploader(ss).Upload(&object)
	if err != nil {
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

func getAWSCredentials(file string) (*credentials.Credentials, error) {
	awsKeyID := ""
	awsSecret := ""

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\[default\]\n+aws_access_key_id\s*=\s*(.*?)\n+aws_secret_access_key\s*=\s*(.*)`)
	if match := re.FindSubmatch(bytes); len(match) == 3 {
		awsKeyID = strings.TrimSpace(string(match[1]))
		awsSecret = strings.TrimSpace(string(match[2]))
	} else {
		re = regexp.MustCompile(`\[default\]\n+aws_secret_access_key\s*=\s*(.*?)\n+aws_access_key_id\s*=\s*(.*)`)
		if match := re.FindSubmatch(bytes); len(match) == 3 {
			awsSecret = strings.TrimSpace(string(match[1]))
			awsKeyID = strings.TrimSpace(string(match[2]))
		}
	}

	if awsKeyID == "" {
		return nil, fmt.Errorf("Invalid AWS credentials - missing 'aws_access_key_id'")
	}

	if awsSecret == "" {
		return nil, fmt.Errorf("Invalid AWS credentials - missing 'aws_secret_access_key'")
	}

	return credentials.NewStaticCredentials(awsKeyID, awsSecret, ""), nil
}
