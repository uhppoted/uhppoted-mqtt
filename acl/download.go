package acl

import (
	"encoding/json"
	"fmt"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"net/url"
)

func (a *ACL) Download(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
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

	acl, err := a.fetch("acl:download", uri.String())
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
