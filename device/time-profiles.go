package device

import (
	"fmt"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetTimeProfile(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID  *uint32 `json:"device-id"`
		ProfileID *uint8  `json:"profile-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.ProfileID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing time profile ID", nil), fmt.Errorf("Invalid/missing time profile ID")
	}

	rq := uhppoted.GetTimeProfileRequest{
		DeviceID:  *body.DeviceID,
		ProfileID: *body.ProfileID,
	}

	response, err := impl.GetTimeProfile(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve time profile %v from %d", *body.ProfileID, *body.DeviceID), err), err
	}

	return response, nil
}
