package acl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"io"
	"net/url"
	"strings"
)

func (a *ACL) Upload(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	devices := ctx.Value("devices").([]*uhppote.Device)

	body := struct {
		URL *string `json:"url"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}, fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.URL == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid upload URI",
		}, fmt.Errorf("Missing/invalid upload URI")
	}

	uri, err := url.Parse(*body.URL)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid upload URI",
		}, fmt.Errorf("Invalid upload URL '%v' (%w)", body.URL, err)
	}

	acl, err := api.GetACL(impl.Uhppote, devices)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving ACL",
		}, err
	}

	if acl == nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving card access permissions",
		}, fmt.Errorf("<nil> response to GetCard request")
	}

	for k, l := range acl {
		a.info("acl:upload", fmt.Sprintf("%v  Retrieved %v records", k, len(l)))
	}

	if err = a.upload(uri.String(), acl, devices); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error uploading ACL",
		}, err
	}

	return struct {
		Uploaded string `json:"uploaded"`
	}{
		Uploaded: uri.String(),
	}, nil
}

func (a *ACL) upload(uri string, acl api.ACL, devices []*uhppote.Device) error {
	var w strings.Builder

	err := api.MakeTSV(acl, devices, &w)
	if err != nil {
		return err
	}

	tsv := []byte(w.String())
	files := map[string][]byte{
		"uhppoted.acl": tsv,
	}

	if signature, err := a.sign(tsv); err != nil {
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

	a.info("acl:upload", fmt.Sprintf("tar'd ACL (%v bytes) and signature (%v bytes): %v bytes", len(files["uhppoted.acl"]), len(files["signature"]), b.Len()))

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

	a.info("acl:upload", fmt.Sprintf("INFO  Stored ACL to %v", uri))

	return nil
}

func (a *ACL) storeHTTP(url string, r io.Reader) error {
	return storeHTTP(url, r)
}

func (a *ACL) storeS3(uri string, r io.Reader) error {
	return storeS3(uri, a.Credentials, a.Region, r)
}

func (a *ACL) storeFile(url string, r io.Reader) error {
	return storeFile(url, r)
}
