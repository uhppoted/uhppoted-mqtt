package device

import (
	"encoding/json"
	"fmt"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetStatus(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
			Debug:   err,
		}, err
	}

	if body.DeviceID == nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Invalid/missing device ID",
		}, fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetStatus(rq)
	if err != nil {
		return common.Error{
			Code:    uhppoted.StatusInternalServerError,
			Message: fmt.Sprintf("Could not retrieve status for %d", *body.DeviceID),
			Debug:   err,
		}, err
	}

	return response, nil
}
