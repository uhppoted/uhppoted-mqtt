package device

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetTime(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Invalid/missing device ID",
		}, fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetTimeRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetTime(rq)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: fmt.Sprintf("Could not retrieve device time for %d", *body.DeviceID),
			Debug:   err,
		}, err
	}

	return response, nil
}

func (d *Device) SetTime(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		DateTime *types.DateTime    `json:"date-time"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Invalid/missing device ID",
		}, fmt.Errorf("Invalid/missing device ID")
	}

	if body.DateTime == nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Invalid/missing datetime",
		}, fmt.Errorf("Invalid/missing datetime")
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: *body.DeviceID,
		DateTime: *body.DateTime,
	}

	response, err := impl.SetTime(rq)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: fmt.Sprintf("Could not set device time for %d", *body.DeviceID),
			Debug:   err,
		}, err
	}

	if response == nil {
		return nil, nil
	}

	return response, nil
}
