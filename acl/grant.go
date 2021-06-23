package acl

import (
	"encoding/json"
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Grant(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32     `json:"card-number"`
		From       *types.Date `json:"start-date"`
		To         *types.Date `json:"end-date"`
		Profile    int         `json:"profile"`
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

	if body.Profile != 0 && (body.Profile < 2 || body.Profile > 254) {
		return common.MakeError(StatusBadRequest, fmt.Sprintf("Invalid time profile (%v)", body.Profile), nil), fmt.Errorf("Invalid time profile (%v)", body.Profile)
	}

	err := api.Grant(impl.UHPPOTE, a.Devices, *body.CardNumber, *body.From, *body.To, body.Profile, body.Doors)
	if err != nil {
		return common.MakeError(StatusInternalServerError, err.Error(), nil), err
	}

	return struct {
		Granted bool `json:"granted"`
	}{
		Granted: true,
	}, nil
}
