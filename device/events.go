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

func (d *Device) GetEvents(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		DeviceID uint32 `json:"device-id"`
		Count    int    `json:"count,omitempty"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
	}

	deviceID := body.DeviceID
	count := body.Count
	events := []any{}

	if count > 0 {
		list, err := impl.GetEvents(deviceID, count)
		if err != nil {
			return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", deviceID), err), err
		}

		for _, e := range list {
			events = append(events, Transmogrify(e))
		}
	}

	first, last, current, err := impl.GetEventIndices(deviceID)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", deviceID), err), err
	}

	response := struct {
		DeviceID uint32 `json:"device-id,omitempty"`
		First    uint32 `json:"first,omitempty"`
		Last     uint32 `json:"last,omitempty"`
		Current  uint32 `json:"current,omitempty"`
		Events   []any  `json:"events,omitempty"`
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
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
	} else {
		deviceID = body.DeviceID
	}

	if body.Index == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing event index", nil), fmt.Errorf("invalid/missing event index")
	}

	// ... parse event index

	if matches := regexp.MustCompile("^([0-9]+|first|last|current|next)$").FindStringSubmatch(fmt.Sprintf("%v", body.Index)); matches == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing event index", nil), fmt.Errorf("invalid/missing event index")
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
			return common.MakeError(StatusBadRequest, fmt.Sprintf("Invalid event index (%v)", body.Index), nil), fmt.Errorf("invalid event index (%v)", index)
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
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("invalid/missing device ID")
	}

	if body.Enabled == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing 'enabled'", nil), fmt.Errorf("invalid/missing 'enable'")
	}

	deviceID := uint32(*body.DeviceID)
	enabled := *body.Enabled

	updated, err := impl.RecordSpecialEvents(deviceID, enabled)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not update 'record special events' flag for %d", *body.DeviceID), err), err
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Enabled  bool   `json:"enabled"`
		Updated  bool   `json:"updated"`
	}{
		DeviceID: deviceID,
		Enabled:  enabled,
		Updated:  updated,
	}

	return response, nil
}

func getEvent(impl uhppoted.IUHPPOTED, deviceID uint32, index uint32) (any, error) {
	event, err := impl.GetEvent(deviceID, index)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve event %v from %v", index, deviceID), err), err
	} else if event == nil {
		return common.MakeError(StatusNotFound, fmt.Sprintf("No event at %v on %v", index, deviceID), nil), fmt.Errorf("no event at %v on %v", index, deviceID)
	}

	response := struct {
		DeviceID uint32 `json:"device-id"`
		Event    any    `json:"event"`
	}{
		DeviceID: deviceID,
		Event:    Transmogrify(*event),
	}

	return &response, nil
}

func getNextEvent(impl uhppoted.IUHPPOTED, deviceID uint32) (any, error) {
	response := struct {
		DeviceID uint32 `json:"device-id"`
		Event    any    `json:"event"`
	}{
		DeviceID: deviceID,
	}

	events, err := impl.GetEvents(deviceID, 1)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve event from %v", deviceID), err), err
	} else if events == nil {
		return common.MakeError(StatusNotFound, fmt.Sprintf("No 'next' event for %v", deviceID), nil), fmt.Errorf("no 'next' event for %v", deviceID)
	} else if len(events) > 0 {
		response.Event = Transmogrify(events[0])
	}

	return &response, nil
}

// FIXME remove when uhppoted-lib::Event is removed
func Transmogrify(e uhppoted.Event) any {
	lookup := func(key string) string {
		if v, ok := locales.Lookup(key); ok {
			return v
		}

		return ""
	}

	return struct {
		DeviceID      uint32         `json:"device-id"`
		Index         uint32         `json:"event-id"`
		Type          uint8          `json:"event-type"`
		TypeText      string         `json:"event-type-text"`
		Granted       bool           `json:"access-granted"`
		Door          uint8          `json:"door-id"`
		Direction     uint8          `json:"direction"`
		DirectionText string         `json:"direction-text"`
		CardNumber    uint32         `json:"card-number"`
		Timestamp     types.DateTime `json:"timestamp"`
		Reason        uint8          `json:"event-reason"`
		ReasonText    string         `json:"event-reason-text"`
	}{
		DeviceID:      e.DeviceID,
		Index:         e.Index,
		Type:          e.Type,
		TypeText:      lookup(fmt.Sprintf("event.type.%v", e.Type)),
		Granted:       e.Granted,
		Door:          e.Door,
		Direction:     e.Direction,
		DirectionText: lookup(fmt.Sprintf("event.direction.%v", e.Direction)),
		CardNumber:    e.CardNumber,
		Timestamp:     e.Timestamp,
		Reason:        e.Reason,
		ReasonText:    lookup(fmt.Sprintf("event.reason.%v", e.Reason)),
	}
}

func transmogrify(e event) any {
	lookup := func(key string) string {
		if v, ok := locales.Lookup(key); ok {
			return v
		}

		return ""
	}

	return struct {
		ControllerID  uint32         `json:"device-id"`
		Index         uint32         `json:"event-id"`
		Type          uint8          `json:"event-type"`
		TypeText      string         `json:"event-type-text"`
		Granted       bool           `json:"access-granted"`
		Door          uint8          `json:"door-id"`
		Direction     uint8          `json:"direction"`
		DirectionText string         `json:"direction-text"`
		CardNumber    uint32         `json:"card-number"`
		Timestamp     types.DateTime `json:"timestamp"`
		Reason        uint8          `json:"event-reason"`
		ReasonText    string         `json:"event-reason-text"`
	}{
		ControllerID:  e.Controller,
		Index:         e.Index,
		Type:          e.Type,
		TypeText:      lookup(fmt.Sprintf("event.type.%v", e.Type)),
		Granted:       e.Granted,
		Door:          e.Door,
		Direction:     e.Direction,
		DirectionText: lookup(fmt.Sprintf("event.direction.%v", e.Direction)),
		CardNumber:    e.CardNumber,
		Timestamp:     e.Timestamp,
		Reason:        e.Reason,
		ReasonText:    lookup(fmt.Sprintf("event.reason.%v", e.Reason)),
	}
}
