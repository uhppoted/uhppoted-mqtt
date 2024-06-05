package device

import (
	"fmt"
	"regexp"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetDoorDelay(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorDelay(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve delay for controller %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}

func (d *Device) SetDoorDelay(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
		Delay    *uint8             `json:"delay"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	if body.Delay == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing delay", nil), fmt.Errorf("invalid/missing delay: %v", body.Delay)
	}

	if *body.Delay == 0 || *body.Door > 60 {
		return common.MakeError(StatusBadRequest, "Invalid/missing delay", nil), fmt.Errorf("invalid/missing delay: %v", *body.Delay)
	}

	deviceID := uint32(*body.DeviceID)
	door := *body.Door
	delay := *body.Delay

	if err := impl.SetDoorDelay(deviceID, door, delay); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not set delay for controller %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Door     uint8  `json:"door"`
		Delay    uint8  `json:"delay"`
	}{
		DeviceID: deviceID,
		Door:     door,
		Delay:    delay,
	}

	return response, nil
}

func (d *Device) GetDoorControl(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorControl(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve control state for controller %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	return response, nil
}

func (d *Device) SetDoorControl(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID  `json:"device-id"`
		Door     *uint8              `json:"door"`
		Control  *types.ControlState `json:"control"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", body.Door)
	}

	if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	if body.Control == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing control state", nil), fmt.Errorf("invalid/missing control state: %v", body.Control)
	}

	if *body.Control < 1 || *body.Control > 3 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door control state", nil), fmt.Errorf("invalid/missing control state: %v", *body.Control)
	}

	deviceID := uint32(*body.DeviceID)
	door := *body.Door
	mode := *body.Control

	if err := impl.SetDoorControl(deviceID, door, mode); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not set delay for controller %d, door %d", *body.DeviceID, *body.Door), err), err
	}

	response := struct {
		DeviceID uint32             `json:"device-id"`
		Door     uint8              `json:"door"`
		Control  types.ControlState `json:"control"`
	}{
		DeviceID: deviceID,
		Door:     door,
		Control:  mode,
	}

	return response, nil
}

// Sets the supervisor passcodes for a door.
//
// Each door may be individually assigned up to four passcodes, with valid passcodes
// being in the range [1..999999]. The function uses the first four codes from the
// supplied list and invalid passcodes are set to 0 (no code).
//
// The controller, door and assigned passcodes are included in the response.
func (d *Device) SetDoorPasscodes(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID  *uhppoted.DeviceID `json:"device-id"`
		Door      *uint8             `json:"door"`
		Passcodes []uint32           `json:"passcodes"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Door == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", body.Door)
	} else if *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	deviceID := uint32(*body.DeviceID)
	door := *body.Door
	passcodes := []uint32{}

	for i := 0; i < 4 && i < len(body.Passcodes); i++ {
		passcode := body.Passcodes[i]
		if passcode > 0 && passcode < 1000000 {
			passcodes = append(passcodes, passcode)
		} else {
			passcodes = append(passcodes, 0)
		}
	}

	if err := impl.SetDoorPasscodes(deviceID, door, passcodes...); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not set passcodes for controller %d, door %d", deviceID, door), err), err
	}

	response := struct {
		DeviceID  uint32   `json:"device-id"`
		Door      uint8    `json:"door"`
		Passcodes []uint32 `json:"passcodes"`
	}{
		DeviceID:  deviceID,
		Door:      door,
		Passcodes: passcodes,
	}

	return response, nil
}

func (d *Device) SetInterlock(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID  *uhppoted.DeviceID `json:"device-id"`
		Interlock *types.Interlock   `json:"interlock"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Interlock == nil {
		return common.MakeError(StatusBadRequest, "Missing door interlock", nil), fmt.Errorf("missing door interlock: %v", body.Interlock)
	}

	deviceID := uint32(*body.DeviceID)
	interlock := *body.Interlock

	if interlock != 0 && interlock != 1 && interlock != 2 && interlock != 3 && interlock != 4 && interlock != 8 {
		return common.MakeError(StatusBadRequest, "Invalid door interlock", nil), fmt.Errorf("invalid door interlock: %v", interlock)
	}

	if err := impl.SetInterlock(deviceID, interlock); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not set interlock for controller %v", deviceID), err), err
	}

	response := struct {
		DeviceID  uint32          `json:"device-id"`
		Interlock types.Interlock `json:"interlock"`
	}{
		DeviceID:  deviceID,
		Interlock: interlock,
	}

	return response, nil
}

func (d *Device) SetKeypads(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Keypads  map[uint8]bool     `json:"keypads"`
	}{
		Keypads: map[uint8]bool{},
	}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	deviceID := uint32(*body.DeviceID)
	keypads := map[uint8]bool{
		1: body.Keypads[1],
		2: body.Keypads[2],
		3: body.Keypads[3],
		4: body.Keypads[4],
	}

	if err := impl.ActivateKeypads(deviceID, keypads); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not activate/deactivate keypads for controller %v", deviceID), err), err
	}

	response := struct {
		DeviceID uint32         `json:"device-id"`
		Keypads  map[uint8]bool `json:"keypads"`
	}{
		DeviceID: deviceID,
		Keypads:  keypads,
	}

	return response, nil
}

func (d *Device) OpenDoor(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Card     *uint32            `json:"card-number"`
		Door     *uint8             `json:"door"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	}

	if body.Card == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing card", nil), fmt.Errorf("invalid/missing card")
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return common.MakeError(StatusBadRequest, "Invalid/missing door", nil), fmt.Errorf("invalid/missing door: %v", *body.Door)
	}

	deviceID := uint32(*body.DeviceID)
	card := *body.Card
	door := *body.Door

	if !d.authorized(card) {
		return common.MakeError(StatusUnauthorized, fmt.Sprintf("Card %v is not authorized for door %v", card, door), nil),
			fmt.Errorf("card %v is not authorized for controller %v, door %v", card, deviceID, door)
	}

	if err := validate(impl, deviceID, card, door); err != nil {
		return common.MakeError(StatusUnauthorized, fmt.Sprintf("Card %v is not authorized for door %v", card, door), nil),
			fmt.Errorf("failed to validate access for card %v to controller %v, door %v (%v)", card, deviceID, door, err)
	}

	rq := uhppoted.OpenDoorRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.OpenDoor(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not open controller %d, door %d with card %v", *body.DeviceID, *body.Door, *body.Card), err), err
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

func validate(impl uhppoted.IUHPPOTED, deviceID uint32, cardNumber uint32, door uint8) error {
	rq := uhppoted.GetCardRequest{
		DeviceID:   uhppoted.DeviceID(deviceID),
		CardNumber: cardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return err
	} else if response == nil {
		return fmt.Errorf("GetCard returned <nil> for card %v, controller %v", cardNumber, deviceID)
	}

	card := response.Card

	// Check start/end validity dates
	today := types.Date(time.Now())
	if card.From.IsZero() || card.To.IsZero() || today.Before(card.From) || today.After(card.To) {
		return fmt.Errorf("card %v is not valid for %v", card.CardNumber, today)
	}

	// Check door permissions
	if card.Doors[door] < 1 || card.Doors[door] > 254 {
		return fmt.Errorf("card %v is does not have permission for %v, door %v", cardNumber, deviceID, door)
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

func getTimeProfile(impl uhppoted.IUHPPOTED, deviceID uint32, profileID uint8) (*types.TimeProfile, error) {
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

func checkTimeProfile(deviceID, cardNumber uint32, profileID uint8, profile types.TimeProfile) error {
	now := time.Now()
	hhmm := types.HHmmFromTime(now)
	today := types.Date(now)

	if profile.From.IsZero() || profile.To.IsZero() || today.Before(profile.From) || today.After(profile.To) {
		return fmt.Errorf("card %v: time profile %v on controller %v is not valid for %v", cardNumber, profileID, deviceID, today)
	}

	if !profile.Weekdays[today.Weekday()] {
		return fmt.Errorf("card %v: time profile %v on controller %v is not authorized for %v", cardNumber, profileID, deviceID, today.Weekday())
	}

	for _, p := range []uint8{1, 2, 3} {
		if segment, ok := profile.Segments[p]; ok {
			if !segment.Start.After(hhmm) && !segment.End.Before(hhmm) {
				return nil
			}
		}
	}

	return fmt.Errorf("card %v: time profile %v on controller %v is not authorized for %v", cardNumber, profileID, deviceID, hhmm)
}
