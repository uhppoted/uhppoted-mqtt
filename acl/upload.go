package acl

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
)

func (a *ACL) Upload(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
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

	acl, err := api.GetACL(impl.Uhppote, a.Devices)
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

	var w strings.Builder
	if err := api.MakeTSV(acl, a.Devices, &w); err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error reformatting card access permissions",
		}, err
	}

	if err = a.store("acl:upload", uri.String(), "uhppoted.acl", []byte(w.String())); err != nil {
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
