package acl

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

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

func unzip(r io.Reader) (map[string][]byte, string, error) {
	files := map[string][]byte{}
	uname := ""

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, "", err
	}

	for _, f := range zr.File {
		if filepath.Ext(f.Name) == ".acl" {
			if _, ok := files["ACL"]; ok {
				return nil, "", fmt.Errorf("Multiple ACL files in tar.gz")
			}

			rc, err := f.Open()
			if err != nil {
				return nil, "", err
			}

			var buffer bytes.Buffer
			if _, err := io.Copy(&buffer, rc); err != nil {
				return nil, "", err
			}

			files["ACL"] = buffer.Bytes()
			uname = f.Comment
			rc.Close()
		}

		if f.Name == "signature" {
			if _, ok := files["signature"]; ok {
				return nil, "", fmt.Errorf("Multiple signature files in tar.gz")
			}

			rc, err := f.Open()
			if err != nil {
				return nil, "", err
			}

			var buffer bytes.Buffer
			if _, err := io.Copy(&buffer, rc); err != nil {
				return nil, "", err
			}

			files["signature"] = buffer.Bytes()
			rc.Close()
		}
	}

	if _, ok := files["ACL"]; !ok {
		return nil, "", fmt.Errorf("ACL file missing from tar.gz")
	}

	if _, ok := files["signature"]; !ok {
		return nil, "", fmt.Errorf("'signature' file missing from tar.gz")
	}

	return files, uname, nil
}
