package acl

import (
	"encoding/json"
	"fmt"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
)

func (a *ACL) Revoke(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32  `json:"card-number"`
		Doors      []string `json:"doors"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}, fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.CardNumber == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}, fmt.Errorf("Missing/invalid card number")
	}

	err := api.Revoke(impl.Uhppote, a.Devices, *body.CardNumber, body.Doors)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: err.Error(),
		}, err
	}

	return struct {
		Revoked bool `json:"revoked"`
	}{
		Revoked: true,
	}, nil
}
