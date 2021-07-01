package acl

import (
	"encoding/json"
	"fmt"

	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Revoke(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32  `json:"card-number"`
		Doors      []string `json:"doors"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.CardNumber == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid card number", nil), fmt.Errorf("Missing/invalid card number")
	}

	err := api.Revoke(a.UHPPOTE, a.Devices, *body.CardNumber, body.Doors)
	if err != nil {
		return common.MakeError(StatusInternalServerError, err.Error(), nil), err
	}

	return struct {
		Revoked bool `json:"revoked"`
	}{
		Revoked: true,
	}, nil
}
