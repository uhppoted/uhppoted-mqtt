package device

import (
	"encoding/json"
	"log"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

const (
	StatusInternalServerError = uhppoted.StatusInternalServerError
	StatusBadRequest          = uhppoted.StatusBadRequest
)

type Device struct {
	Log *log.Logger
}

func unmarshal(bytes []byte, request interface{}) (interface{}, error) {
	if err := json.Unmarshal(bytes, request); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), err
	}

	return nil, nil
}
