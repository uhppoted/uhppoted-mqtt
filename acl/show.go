package acl

import (
	"encoding/json"
	"fmt"

	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Show(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), fmt.Errorf("%w: %v", uhppoted.ErrBadRequest, err)
	}

	if body.CardNumber == nil {
		return common.MakeError(StatusBadRequest, "Missing/invalid card number", nil), fmt.Errorf("Missing/invalid card number")
	}

	acl, err := api.GetCard(a.UHPPOTE, a.Devices, *body.CardNumber)
	if err != nil {
		return common.MakeError(StatusInternalServerError, "Error retrieving card access permissions", err), err
	}

	if acl == nil {
		return common.MakeError(StatusInternalServerError, "Error retrieving card access permissions", nil), fmt.Errorf("<nil> response to GetCard request")
	}

	response := Permissions{
		CardNumber:  *body.CardNumber,
		Permissions: []Permission{},
	}

	for k, v := range acl {
		response.Permissions = append(response.Permissions, Permission{
			Door:      k,
			StartDate: v.From,
			EndDate:   v.To,
			Profile:   v.Profile,
		})
	}

	return response, nil
}
