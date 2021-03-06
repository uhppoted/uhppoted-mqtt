package acl

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Upload(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		URL *string `json:"url"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.URL == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid upload URI", nil), fmt.Errorf("Missing/invalid upload URI")
	}

	uri, err := url.Parse(*body.URL)
	if err != nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid upload URI", err), fmt.Errorf("Invalid upload URL '%v' (%w)", body.URL, err)
	}

	acl, err := api.GetACL(impl.Uhppote, a.Devices)
	if err != nil {
		return common.MakeError(StatusInternalServerError, "Error retrieving ACL", err), err
	}

	if acl == nil {
		return common.MakeError(StatusInternalServerError, "Error retrieving card access permissions", nil), fmt.Errorf("<nil> response to GetCard request")
	}

	for k, l := range acl {
		a.info("acl:upload", fmt.Sprintf("%v  Retrieved %v records", k, len(l)))
	}

	var w strings.Builder
	if err := api.MakeTSV(acl, a.Devices, &w); err != nil {
		return common.MakeError(StatusInternalServerError, "Error reformatting card access permissions", err), err
	}

	if err = a.store("acl:upload", uri.String(), "uhppoted.acl", []byte(w.String())); err != nil {
		return common.MakeError(StatusBadRequest, "Error uploading ACL", err), err
	}

	return struct {
		Uploaded string `json:"uploaded"`
	}{
		Uploaded: uri.String(),
	}, nil
}
