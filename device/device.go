package device

import (
	"encoding/json"
	"log"

	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

const (
	StatusInternalServerError = uhppoted.StatusInternalServerError
	StatusBadRequest          = uhppoted.StatusBadRequest
	StatusUnauthorized        = uhppoted.StatusUnauthorized
	StatusNotFound            = uhppoted.StatusNotFound
)

type Device struct {
	AuthorizedCards []string
	Log             *log.Logger
}

var protocol string = ""

func SetProtocol(version string) {
	protocol = version
}

func unmarshal(bytes []byte, request interface{}) (interface{}, error) {
	if err := json.Unmarshal(bytes, request); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), err
	}

	return nil, nil
}
