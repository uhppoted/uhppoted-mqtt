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
		return common.MakeError(StatusInternalServerError, "Error searching for active devices", err), err
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
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetDevice(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve device information for %d", *body.DeviceID), err), err
	}

	return response, nil
}
