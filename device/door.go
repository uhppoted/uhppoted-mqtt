package device

import (
	"fmt"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetDoorDelay(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", *body.Door)
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorDelay(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve delay for device %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}

func (d *Device) SetDoorDelay(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
		Delay    *uint8             `json:"delay"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", *body.Door)
	}

	if body.Delay == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing delay", nil), fmt.Errorf("Invalid/missing delay: %v", body.Delay)
	}

	if *body.Delay == 0 || *body.Door > 60 {
		return common.MakeError(StatusBadRequest, "Invalid/missing delay", nil), fmt.Errorf("Invalid/missing delay: %v", *body.Delay)
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Delay:    *body.Delay,
	}

	response, err := impl.SetDoorDelay(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not setting delay for device %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}

func (d *Device) GetDoorControl(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", *body.Door)
	}

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorControl(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve control state for device %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}

func (d *Device) SetDoorControl(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID     `json:"device-id"`
		Door     *uint8                 `json:"door"`
		Control  *uhppoted.ControlState `json:"control"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", *body.Door)
	}

	if body.Control == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing control state", nil), fmt.Errorf("Invalid/missing control state: %v", body.Control)
	}

	if *body.Control < 1 || *body.Control > 3 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door control state", nil), fmt.Errorf("Invalid/missing control state: %v", *body.Control)
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Control:  *body.Control,
	}

	response, err := impl.SetDoorControl(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not setting delay for device %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}
