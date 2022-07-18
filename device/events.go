package device

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/locales"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

type Event struct {
	DeviceID uint32 `json:"device-id"`
	Index    uint32 `json:"event-id"`
	Type     struct {
		Code        uint8  `json:"code"`
		Description string `json:"description"`
	} `json:"event-type"`
	Granted   bool  `json:"access-granted"`
	Door      uint8 `json:"door-id"`
	Direction struct {
		Code        uint8  `json:"code"`
		Description string `json:"description"`
	} `json:"direction"`
	CardNumber uint32         `json:"card-number"`
	Timestamp  types.DateTime `json:"timestamp"`
	Reason     struct {
		Code        uint8  `json:"code"`
		Description string `json:"description"`
	} `json:"event-reason"`
}

func (d *Device) GetEvents(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID uint32 `json:"device-id"`
		Count    int    `json:"count,omitempty"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	deviceID := body.DeviceID
	count := body.Count
	events := []Event{}

	if count > 0 {
		list, err := impl.GetEvents(deviceID, count)
		if err != nil {
			return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", deviceID), err), err
		}

		for _, e := range list {
			events = append(events, transmogrify(e))
		}
	}

	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", deviceID), err), err
	}

	response := struct {
		DeviceID uint32  `json:"device-id,omitempty"`
		First    uint32  `json:"first,omitempty"`
		Last     uint32  `json:"last,omitempty"`
		Current  uint32  `json:"current,omitempty"`
		Events   []Event `json:"events,omitempty"`
	}{
		DeviceID: deviceID,
		First:    first,
		Last:     last,
		Current:  current,
		Events:   events,
	}

	return response, nil
}

func (d *Device) GetEvent(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	var deviceID uint32
	var index string

	body := struct {
		DeviceID uint32 `json:"device-id"`
		Index    any    `json:"event-index"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	} else {
		deviceID = body.DeviceID
	}

	if body.Index == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing event index", nil), fmt.Errorf("Invalid/missing event index")
	}

	// ... parse event index

	if matches := regexp.MustCompile("^([0-9]+|first|last|current|next)$").FindStringSubmatch(fmt.Sprintf("%v", body.Index)); matches == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing event index", nil), fmt.Errorf("Invalid/missing event index")
	} else {
		index = matches[1]
	}

	// .. get event indices
	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve event indices from %v", deviceID), err), err
	}

	// ... get event
	switch index {
	case "first":
		return getEvent(impl, deviceID, first)

	case "last":
		return getEvent(impl, deviceID, last)

	case "current":
		return getEvent(impl, deviceID, current)

	case "next":
		return getNextEvent(impl, deviceID)

	default:
		if v, err := strconv.ParseUint(index, 10, 32); err != nil {
			return common.MakeError(StatusBadRequest, fmt.Sprintf("Invalid event index (%v)", body.Index), nil), fmt.Errorf("Invalid event index (%v)", index)
		} else {
			return getEvent(impl, deviceID, uint32(v))
		}
	}
}

// Handler for the special-events MQTT message. Extracts the 'enabled' value from the request
// and invokes the uhppoted-lib.RecordSpecialEvents API function to update the controller
// 'record special events' flag.
func (d *Device) RecordSpecialEvents(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Enabled  *bool              `json:"enabled"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Enabled == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing 'enabled'", nil), fmt.Errorf("Invalid/missing 'enabled'")
	}

	rq := uhppoted.RecordSpecialEventsRequest{
		DeviceID: *body.DeviceID,
		Enable:   *body.Enabled,
	}

	response, err := impl.RecordSpecialEvents(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not update 'record special events' flag for %d", *body.DeviceID), err), err
	}

	if response == nil {
		return nil, nil
	}

	return response, nil
}

func getEvent(impl uhppoted.IUHPPOTED, deviceID uint32, index uint32) (any, error) {
	event, err := impl.GetEvent(deviceID, index)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve event %v from %v", index, deviceID), err), err
	} else if event == nil {
		return common.MakeError(StatusNotFound, fmt.Sprintf("No event at %v on %v", index, deviceID), nil), fmt.Errorf("No event at %v on %v", index, deviceID)
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Event    any    `json:"event"`
	}{
		DeviceID: deviceID,
		Event:    transmogrify(*event),
	}

	return &response, nil
}

func getNextEvent(impl uhppoted.IUHPPOTED, deviceID uint32) (any, error) {
	event, err := impl.GetEvent(deviceID, 1)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve event from %v", deviceID), err), err
	} else if event == nil {
		return common.MakeError(StatusNotFound, fmt.Sprintf("No 'next' event for %v", deviceID), nil), fmt.Errorf("No 'next' event for %v", deviceID)
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Event    any    `json:"event"`
	}{
		DeviceID: deviceID,
		Event:    transmogrify(*event),
	}

	return &response, nil
}

func transmogrify(e uhppoted.Event) Event {
	return Event{
		DeviceID: e.DeviceID,
		Index:    e.Index,
		Type: struct {
			Code        uint8  `json:"code"`
			Description string `json:"description"`
		}{
			Code:        e.Type,
			Description: lookup(fmt.Sprintf("event.type.%v", e.Type)),
		},
		Granted: e.Granted,
		Door:    e.Door,
		Direction: struct {
			Code        uint8  `json:"code"`
			Description string `json:"description"`
		}{
			Code:        e.Direction,
			Description: lookup(fmt.Sprintf("event.direction.%v", e.Direction)),
		},
		CardNumber: e.CardNumber,
		Timestamp:  e.Timestamp,
		Reason: struct {
			Code        uint8  `json:"code"`
			Description string `json:"description"`
		}{
			Code:        e.Reason,
			Description: lookup(fmt.Sprintf("event.reason.%v", e.Reason)),
		},
	}
}

func lookup(key string) string {
	if v, ok := locales.Lookup(key); ok {
		return v
	}

	return ""
}
