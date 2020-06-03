package acl

import (
	"encoding/json"
	"fmt"

	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (a *ACL) Show(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}, fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.CardNumber == nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}, fmt.Errorf("Missing/invalid card number")
	}

	acl, err := api.GetCard(impl.Uhppote, a.Devices, *body.CardNumber)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving card access permissions",
		}, err
	}

	if acl == nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving card access permissions",
		}, fmt.Errorf("<nil> response to GetCard request")
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
		})
	}

	return response, nil
}
