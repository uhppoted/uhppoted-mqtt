package acl

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

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

func untar(r io.Reader) (map[string][]byte, string, error) {
	files := map[string][]byte{}
	uname := ""

	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, "", err
	}

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, "", err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			if filepath.Ext(header.Name) == ".acl" {
				if _, ok := files["ACL"]; ok {
					return nil, "", fmt.Errorf("multiple ACL files in tar.gz")
				}

				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, tr); err != nil {
					return nil, "", err
				}

				files["ACL"] = buffer.Bytes()
				uname = header.Uname
			}

			if header.Name == "signature" {
				if _, ok := files["signature"]; ok {
					return nil, "", fmt.Errorf("multiple signature files in tar.gz")
				}

				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, tr); err != nil {
					return nil, "", err
				}

				files["signature"] = buffer.Bytes()
			}
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
