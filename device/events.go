package device

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

type startdate time.Time
type enddate time.Time

func (d *Device) GetEvents(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Start    *startdate         `json:"start"`
		End      *enddate           `json:"end"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Start != nil && body.End != nil && time.Time(*body.End).Before(time.Time(*body.Start)) {
		return common.MakeError(StatusBadRequest, "Invalid event data range", nil), fmt.Errorf("Invalid event date range: %v to %v", body.Start, body.End)
	}

	rq := uhppoted.GetEventRangeRequest{
		DeviceID: *body.DeviceID,
		Start:    (*types.DateTime)(body.Start),
		End:      (*types.DateTime)(body.End),
	}

	response, err := impl.GetEventRange(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) GetEvent(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		EventID  *uint32            `json:"event-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.EventID == nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing event ID", nil), fmt.Errorf("Invalid/missing event ID")
	}

	rq := uhppoted.GetEventRequest{
		DeviceID: *body.DeviceID,
		EventID:  *body.EventID,
	}

	response, err := impl.GetEvent(rq)
	if err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve events from %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *startdate) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		*d = startdate(datetime)
		return nil
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		*d = startdate(datetime)
		return nil
	}

	if date, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		*d = startdate(date)
		return nil
	}

	return fmt.Errorf("Cannot parse date/time %s", string(bytes))
}

func (d *enddate) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		*d = enddate(datetime)
		return nil
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		*d = enddate(datetime)
		return nil
	}

	if date, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		*d = enddate(time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.Local))
		return nil
	}

	return fmt.Errorf("Cannot parse date/time %s", string(bytes))
}
