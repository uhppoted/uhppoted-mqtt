package device

import (
	"encoding/json"

	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
	"github.com/uhppoted/uhppoted-mqtt/log"
)

const (
	StatusInternalServerError = uhppoted.StatusInternalServerError
	StatusBadRequest          = uhppoted.StatusBadRequest
	StatusUnauthorized        = uhppoted.StatusUnauthorized
	StatusNotFound            = uhppoted.StatusNotFound
)

type Device struct {
	AuthorizedCards []string
}

func SetProtocol(version string) {
}

func unmarshal(bytes []byte, request interface{}) (interface{}, error) {
	if err := json.Unmarshal(bytes, request); err != nil {
		return common.MakeError(StatusBadRequest, "Cannot parse request", err), err
	}

	return nil, nil
}

func infof(format string, args ...any) {
	log.Infof("mqttd", format, args...)
}

func warnf(format string, args ...any) {
	log.Warnf("mqttd", format, args...)
}
