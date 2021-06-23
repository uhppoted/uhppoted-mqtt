package device

import (
	"fmt"
	"regexp"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
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

func (d *Device) OpenDoor(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Card     *uint32            `json:"card-number"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Card == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing card", nil), fmt.Errorf("Invalid/missing card")
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("Invalid/missing door: %v", *body.Door)
	}

	deviceID := uint32(*body.DeviceID)
	card := *body.Card
	door := *body.Door

	if !d.authorized(card) {
		return common.MakeError(StatusUnauthorized, fmt.Sprintf("Card %v is not authorized for door %v", card, door), nil),
			fmt.Errorf("Card %v is not authorized for device %v, door %v", card, deviceID, door)
	}

	if err := validate(impl, deviceID, card, door); err != nil {
		return common.MakeError(StatusUnauthorized, fmt.Sprintf("Card %v is not authorized for door %v", card, door), nil),
			fmt.Errorf("Failed to validate access for card %v to device %v, door %v (%v)", card, deviceID, door, err)
	}

	rq := uhppoted.OpenDoorRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.OpenDoor(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not open device %d, door %d with card %v", *body.DeviceID, *body.Door, *body.Card), err), err
	}

	return response, nil
}

func (d *Device) authorized(card uint32) bool {
	c := fmt.Sprintf("%v", card)
	for _, re := range d.AuthorizedCards {
		if ok, err := regexp.MatchString(re, c); ok && err == nil {
			return true
		}
	}

	return false
}

func validate(impl *uhppoted.UHPPOTED, deviceID uint32, cardNumber uint32, door uint8) error {
	rq := uhppoted.GetCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return err
	} else if response == nil {
		return fmt.Errorf("GetCard returned <nil> for card %v, device %v", cardNumber, deviceID)
	}

	card := response.Card

	// Check start/end validity dates
	today := types.Date(time.Now())
	if card.From == nil || card.To == nil || today.Before(*card.From) || today.After(*card.To) {
		return fmt.Errorf("Card %v is not valid for %v", card.CardNumber, today)
	}

	// Check door permissions
	if card.Doors[door] < 1 || card.Doors[door] > 254 {
		return fmt.Errorf("Card %v is does not have permission for %v, door %v", cardNumber, deviceID, door)
	}

	// Check time profile
	if card.Doors[door] >= 2 && card.Doors[door] <= 254 {
		profileID := uint8(card.Doors[door])
		checked := map[uint8]bool{}

		for {
			profile, err := getTimeProfile(impl, deviceID, profileID)
			if err != nil {
				return err
			}

			if profile == nil {
				return fmt.Errorf("GetTimeProfile received <nil> response for time profile %v associated with card %v, door %v from device %v", profileID, cardNumber, door, deviceID)
			}

			if err = checkTimeProfile(deviceID, cardNumber, card.Doors[door], *profile); err == nil {
				break
			}

			if profile.LinkedProfileID < 2 || profile.LinkedProfileID > 254 || checked[profile.LinkedProfileID] {
				return err
			}

			checked[profileID] = true
			profileID = profile.LinkedProfileID
		}
	}

	return nil
}

func getTimeProfile(impl *uhppoted.UHPPOTED, deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
	rq := uhppoted.GetTimeProfileRequest{
		DeviceID:  deviceID,
		ProfileID: profileID,
	}

	response, err := impl.GetTimeProfile(rq)
	if err != nil {
		return nil, err
	}

	return &response.TimeProfile, nil
}

func checkTimeProfile(deviceID, cardNumber uint32, profileID int, profile types.TimeProfile) error {
	now := time.Now()
	hhmm := types.HHmmFromTime(now)
	today := types.Date(now)

	if profile.From == nil || profile.To == nil || today.Before(*profile.From) || today.After(*profile.To) {
		return fmt.Errorf("Card %v: time profile %v on device %v is not valid for %v", cardNumber, profileID, deviceID, today)
	}

	if !profile.Weekdays[today.Weekday()] {
		return fmt.Errorf("Card %v: time profile %v on device %v is not authorized for %v", cardNumber, profileID, deviceID, today.Weekday())
	}

	for _, p := range []uint8{1, 2, 3} {
		if segment, ok := profile.Segments[p]; ok {
			if !segment.Start.After(hhmm) && !segment.End.Before(hhmm) {
				return nil
			}
		}
	}

	return fmt.Errorf("Card %v: time profile %v on device %v is not authorized for %v", cardNumber, profileID, deviceID, hhmm)
}
