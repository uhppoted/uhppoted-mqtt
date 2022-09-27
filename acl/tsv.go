package acl

import (
	"io"
)

func unpackTSV(r io.Reader) (map[string][]byte, string, error) {
	files := map[string][]byte{}
	uname := ""

	if bytes, err := io.ReadAll(r); err != nil {
		return nil, "", err
	} else {
		files["ACL"] = bytes
	}

	files["signature"] = []byte{}

	return files, uname, nil
}
