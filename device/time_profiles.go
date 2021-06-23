package device

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetTimeProfile(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
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

func (d *Device) PutTimeProfile(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uint32            `json:"device-id"`
		Profile  *types.TimeProfile `json:"profile"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Profile == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing time profile", nil), fmt.Errorf("Invalid/missing time profile")
	}

	rq := uhppoted.PutTimeProfileRequest{
		DeviceID:    *body.DeviceID,
		TimeProfile: *body.Profile,
	}

	response, err := impl.PutTimeProfile(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not store time profile %v to %d", body.Profile.ID, *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) ClearTimeProfiles(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.ClearTimeProfilesRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.ClearTimeProfiles(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not clear time profiles from %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) GetTimeProfiles(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		From     int     `json:"from"`
		To       int     `json:"to"`
	}{
		From: 2,
		To:   254,
	}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetTimeProfilesRequest{
		DeviceID: *body.DeviceID,
		From:     body.From,
		To:       body.To,
	}

	response, err := impl.GetTimeProfiles(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve time profiles from %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) PutTimeProfiles(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uint32             `json:"device-id"`
		Profiles []types.TimeProfile `json:"profiles"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.PutTimeProfilesRequest{
		DeviceID: *body.DeviceID,
		Profiles: body.Profiles,
	}

	response, _, err := impl.PutTimeProfiles(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not store time profile set to %d", *body.DeviceID), err), err
	}

	warnings := []string{}
	for _, w := range response.Warnings {
		warnings = append(warnings, fmt.Sprintf("%v", w))
	}

	return struct {
		DeviceID uint32   `json:"device-id"`
		Warnings []string `json:"warnings"`
	}{
		DeviceID: uint32(response.DeviceID),
		Warnings: warnings,
	}, nil
}
