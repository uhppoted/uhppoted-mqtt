package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

type Permission struct {
	Door      string     `json:"door"`
	StartDate types.Date `json:"start-date"`
	EndDate   types.Date `json:"end-date"`
}

type Permissions struct {
	CardNumber  uint32       `json:"card-number"`
	Permissions []Permission `json:"permissions"`
}

func Show(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	meta := ctx.Value("metainfo").(common.MetaInfo)
	devices := ctx.Value("devices").([]*uhppote.Device)

	body := struct {
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return struct {
			Code    int    `json:"error-code"`
			Message string `json:"message"`
		}{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}, fmt.Errorf("%w: %v", uhppoted.BadRequest, err)
	}

	if body.CardNumber == nil {
		return struct {
			Code    int    `json:"error-code"`
			Message string `json:"message"`
		}{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid card number",
		}, fmt.Errorf("Missing/invalid card number")
	}

	acl, err := api.GetCard(impl.Uhppote, devices, *body.CardNumber)
	if err != nil {
		return struct {
			Code    int    `json:"error-code"`
			Message string `json:"message"`
		}{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error retrieving card access permissions",
		}, err
	}

	if acl == nil {
		return struct {
			Code    int    `json:"error-code"`
			Message string `json:"message"`
		}{
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

	return struct {
		common.MetaInfo
		Permissions
	}{
		MetaInfo:    meta,
		Permissions: response,
	}, nil
}
