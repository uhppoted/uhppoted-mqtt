package device

import (
	"fmt"

	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetStatus(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetStatus(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve status for %d", *body.DeviceID), err), err
	}

	return response, nil
}
