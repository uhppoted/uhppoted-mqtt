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
	"log"
	"net/url"
	"strings"
)

type uploader struct {
	keyfile     string
	credentials string
	region      string
}

func Upload(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	devices := ctx.Value("devices").([]*uhppote.Device)
	log := ctx.Value("log").(*log.Logger)

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
		}, fmt.Errorf("Invalid upload URL '%s' (%w)", body.URL, err)
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
		log.Printf("INFO  %v  Retrieved %v records\n", k, len(l))
	}

	u := uploader{
		keyfile:     DEFAULT_KEYFILE,
		credentials: DEFAULT_CREDENTIALS,
		region:      DEFAULT_REGION,
	}

	if err = u.upload(uri.String(), acl, devices, log); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error uploading ACL",
		}, err
	}

	return struct {
		Uploaded bool `json:"uploaded"`
	}{
		Uploaded: true,
	}, nil
}

func (u *uploader) upload(uri string, acl api.ACL, devices []*uhppote.Device, log *log.Logger) error {
	var w strings.Builder

	err := api.MakeTSV(acl, devices, &w)
	if err != nil {
		return err
	}

	tsv := []byte(w.String())
	files := map[string][]byte{
		"uhppoted.acl": tsv,
	}

	if signature, err := sign(tsv, u.keyfile); err != nil {
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

	log.Printf("INFO  tar'd ACL (%v bytes) and signature (%v bytes): %v bytes", len(files["uhppoted.acl"]), len(files["signature"]), b.Len())

	f := u.storeHTTP
	if strings.HasPrefix(uri, "s3://") {
		f = u.storeS3
	} else if strings.HasPrefix(uri, "file://") {
		f = u.storeFile
	}

	if err := f(uri, bytes.NewReader(b.Bytes())); err != nil {
		return err
	}

	log.Printf("INFO  Stored ACL to %v", uri)

	return nil
}

func (u *uploader) storeHTTP(url string, r io.Reader) error {
	return storeHTTP(url, r)
}

func (u *uploader) storeS3(uri string, r io.Reader) error {
	return storeS3(uri, u.credentials, u.region, r)
}

func (u *uploader) storeFile(url string, r io.Reader) error {
	return storeFile(url, r)
}
