package device

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

type Status struct {
	DoorState      map[uint8]bool `json:"door-states"`
	DoorButton     map[uint8]bool `json:"door-buttons"`
	SystemError    uint8          `json:"system-error"`
	SystemDateTime types.DateTime `json:"system-datetime"`
	SequenceId     uint32         `json:"sequence-id"`
	SpecialInfo    uint8          `json:"special-info"`
	RelayState     uint8          `json:"relay-state"`
	InputState     uint8          `json:"input-state"`
	Event          any            `json:"event,omitempty"`
}

func (d *Device) GetStatus(impl uhppoted.IUHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID uint32 `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
	}

	deviceID := body.DeviceID

	reply, err := impl.GetStatus(deviceID)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve status for %d", deviceID), err), err
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Status   Status `json:"status"`
	}{
		DeviceID: deviceID,
		Status: Status{
			DoorState:      map[uint8]bool{},
			DoorButton:     map[uint8]bool{},
			SystemError:    reply.SystemError,
			SystemDateTime: reply.SystemDateTime,
			SequenceId:     reply.SequenceId,
			SpecialInfo:    reply.SpecialInfo,
			RelayState:     reply.RelayState,
			InputState:     reply.InputState,
			Event:          nil,
		},
	}

	for k, v := range reply.DoorState {
		response.Status.DoorState[k] = v
	}

	for k, v := range reply.DoorButton {
		response.Status.DoorButton[k] = v
	}

	if reply.Event != nil {
		event := Transmogrify(*reply.Event)
		response.Status.Event = &event
	}

	return &response, nil
}
