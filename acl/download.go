package acl

import (
	"encoding/json"
	"fmt"
	"net/url"

	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Download(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		URL      *string `json:"url"`
		MimeType string  `json:"mime-type"`
	}{
		MimeType: "application/tar+gzip",
	}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), fmt.Errorf("%w: %v", uhppoted.ErrBadRequest, err)
	}

	if body.URL == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid download URL", nil), fmt.Errorf("missing/invalid download URL")
	}

	uri, err := url.Parse(*body.URL)
	if err != nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid download URL", err), fmt.Errorf("invalid download URL '%v' (%w)", body.URL, err)
	}

	acl, err := a.fetch("acl:download", uri.String(), body.MimeType)
	if err != nil {
		return common.MakeError(StatusBadRequest, "Error downloading ACL", err), err
	}

	if acl == nil {
		return common.MakeError(StatusBadRequest, "Error downloading ACL", nil), fmt.Errorf("download return nil ACL")
	}

	for k, l := range *acl {
		infof("acl:download", "%v  Retrieved %v records", k, len(l))
	}

	rpt, errors := api.PutACL(a.UHPPOTE, *acl, false)
	if len(errors) > 0 {
		err := fmt.Errorf("%v", errors)
		return common.MakeError(StatusInternalServerError, "Error updating ACL", err), err
	}

	summary := map[uint32]struct {
		Unchanged int `json:"unchanged"`
		Updated   int `json:"updated"`
		Added     int `json:"added"`
		Deleted   int `json:"deleted"`
		Failed    int `json:"failed"`
		Errors    int `json:"errors"`
	}{}

	for k, v := range rpt {
		infof("acl:download", "%v  SUMMARY  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v",
			k,
			len(v.Unchanged),
			len(v.Updated),
			len(v.Added),
			len(v.Deleted),
			len(v.Failed),
			len(v.Errors))

		summary[k] = struct {
			Unchanged int `json:"unchanged"`
			Updated   int `json:"updated"`
			Added     int `json:"added"`
			Deleted   int `json:"deleted"`
			Failed    int `json:"failed"`
			Errors    int `json:"errors"`
		}{
			Unchanged: len(v.Unchanged),
			Updated:   len(v.Updated),
			Added:     len(v.Added),
			Deleted:   len(v.Deleted),
			Failed:    len(v.Failed),
			Errors:    len(v.Errors),
		}
	}

	warnings := []string{}
	duplicates := map[string]bool{}
	for k, v := range rpt {
		for _, err := range v.Errors {
			warning := fmt.Sprintf("%v: %v", k, err)
			if _, ok := duplicates[warning]; !ok {
				warnings = append(warnings, warning)
				duplicates[warning] = true
			}
		}
	}

	return struct {
		Report map[uint32]struct {
			Unchanged int `json:"unchanged"`
			Updated   int `json:"updated"`
			Added     int `json:"added"`
			Deleted   int `json:"deleted"`
			Failed    int `json:"failed"`
			Errors    int `json:"errors"`
		} `json:"report"`
		Warnings []string `json:"warnings,omitempty"`
	}{
		Report:   summary,
		Warnings: warnings,
	}, nil

}
