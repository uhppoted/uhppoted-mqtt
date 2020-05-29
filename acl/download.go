package acl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	//	"io"
	"net/url"
	"strings"
)

func (a *ACL) Download(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
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
			Message: "Missing/invalid download URL",
		}, fmt.Errorf("Missing/invalid download URL")
	}

	uri, err := url.Parse(*body.URL)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid download URL",
		}, fmt.Errorf("Invalid download URL '%v' (%w)", body.URL, err)
	}

	acl, err := a.download(uri.String(), devices)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error downloading ACL",
		}, err
	}

	if acl == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Error downloading ACL",
		}, fmt.Errorf("Download return nil ACL")
	}

	for k, l := range *acl {
		a.info("acl:download", fmt.Sprintf("%v  Retrieved %v records", k, len(l)))
	}

	rpt, err := api.PutACL(impl.Uhppote, *acl)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error updating ACL",
		}, err
	}

	summary := map[uint32]struct {
		Unchanged int `json:"unchanged"`
		Updated   int `json:"updated"`
		Added     int `json:"added"`
		Deleted   int `json:"deleted"`
		Failed    int `json:"failed"`
	}{}

	for k, v := range rpt {
		a.info("acl:download", fmt.Sprintf("%v  SUMMARY  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v", k, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed))
		summary[k] = struct {
			Unchanged int `json:"unchanged"`
			Updated   int `json:"updated"`
			Added     int `json:"added"`
			Deleted   int `json:"deleted"`
			Failed    int `json:"failed"`
		}{
			Unchanged: v.Unchanged,
			Updated:   v.Updated,
			Added:     v.Added,
			Deleted:   v.Deleted,
			Failed:    v.Failed,
		}
	}

	return struct {
		Report map[uint32]struct {
			Unchanged int `json:"unchanged"`
			Updated   int `json:"updated"`
			Added     int `json:"added"`
			Deleted   int `json:"deleted"`
			Failed    int `json:"failed"`
		} `json:"report"`
	}{
		Report: summary,
	}, nil

}

func (a *ACL) download(uri string, devices []*uhppote.Device) (*api.ACL, error) {
	a.info("acl:download", fmt.Sprintf("Fetching ACL from %v", uri))

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

	a.info("acl:download", fmt.Sprintf("Fetched ACL from %v (%d bytes)", uri, len(b)))

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

	a.info("acl:download", fmt.Sprintf("Extracted ACL from %v: %v bytes, signature: %v bytes", uri, len(tsv), len(signature)))

	if !a.NoVerify {
		if err := a.verify(uname, tsv, signature); err != nil {
			return nil, err
		}
	}

	acl, err := api.ParseTSV(bytes.NewReader(tsv), devices)
	if err != nil {
		return nil, err
	}

	return &acl, nil
}

func (a *ACL) fetchHTTP(url string) ([]byte, error) {
	return fetchHTTP(url)
}

func (a *ACL) fetchS3(uri string) ([]byte, error) {
	return fetchS3(uri, a.Credentials, a.Region)
}

func (a *ACL) fetchFile(url string) ([]byte, error) {
	return fetchFile(url)
}
