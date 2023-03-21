package device

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) PutTaskList(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uint32      `json:"device-id"`
		Tasks    []types.Task `json:"tasks"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
	}

	rq := uhppoted.PutTaskListRequest{
		DeviceID: *body.DeviceID,
		Tasks:    body.Tasks,
	}

	response, _, err := impl.PutTaskList(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("%d: could not update task list", *body.DeviceID), err), err
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
