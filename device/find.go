package device

import (
	"fmt"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetDevices(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: "Error searching for active devices",
			Debug:   fmt.Errorf("%w: %v", uhppoted.StatusInternalServerError, err),
		}, err
	}

	return response, nil
}

func (d *Device) GetDevice(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
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

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetDevice(rq)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: fmt.Sprintf("Could not retrieve device information for %d", *body.DeviceID),
			Debug:   err,
		}, err
	}

	return response, nil
}
