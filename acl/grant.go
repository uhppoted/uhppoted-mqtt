package acl

import (
	"encoding/json"
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Grant(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32     `json:"card-number"`
		From       *types.Date `json:"start-date"`
		To         *types.Date `json:"end-date"`
		Doors      []string    `json:"doors"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.CardNumber == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid card number", nil), fmt.Errorf("Missing/invalid card number")
	}

	if body.From == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid start date", nil), fmt.Errorf("Missing/invalid start date")
	}

	if body.To == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid end date", nil), fmt.Errorf("Missing/invalid end date")
	}

	err := api.Grant(impl.Uhppote, a.Devices, *body.CardNumber, *body.From, *body.To, body.Doors)
	if err != nil {
		return common.MakeError(StatusInternalServerError, err.Error(), nil), err
	}

	return struct {
		Granted bool `json:"granted"`
	}{
		Granted: true,
	}, nil
}
