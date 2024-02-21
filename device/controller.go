package device

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetDevices(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, "Error searching for active devices", err), err
	}

	return response, nil
}

func (d *Device) GetDevice(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
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

// Resets a controller to the manufacturer default configuration.
func (d *Device) RestoreDefaultParameters(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	controller := uint32(*body.DeviceID)

	if err := impl.RestoreDefaultParameters(controller); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not reset controller %d to manufacturer default configuration", controller), err), err
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Reset    bool   `json:"reset"`
	}{
		DeviceID: controller,
		Reset:    true,
	}

	return response, nil
}
