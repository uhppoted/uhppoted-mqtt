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

func Grant(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	meta := ctx.Value("metainfo").(common.MetaInfo)
	devices := ctx.Value("devices").([]*uhppote.Device)

	body := struct {
		CardNumber *uint32     `json:"card-number"`
		From       *types.Date `json:"start-date"`
		To         *types.Date `json:"end-date"`
		Doors      []string    `json:"doors"`
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

	if body.From == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid start date",
		}, fmt.Errorf("Missing/invalid start date")
	}

	if body.To == nil {
		return Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid end date",
		}, fmt.Errorf("Missing/invalid end date")
	}

	err := api.Grant(impl.Uhppote, devices, *body.CardNumber, *body.From, *body.To, body.Doors)
	if err != nil {
		return Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: err.Error(),
		}, err
	}

	return struct {
		common.MetaInfo
		Granted bool `json:"granted"`
	}{
		MetaInfo: meta,
		Granted:  true,
	}, nil
}
