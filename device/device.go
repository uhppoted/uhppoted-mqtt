package device

import (
	"encoding/json"
	"log"

	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

type Device struct {
	Log *log.Logger
}

func unmarshal(bytes []byte, request interface{}) (interface{}, error) {
	if err := json.Unmarshal(bytes, request); err != nil {
		return common.Error{
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
			Debug:   err,
		}, err
	}

	return nil, nil
}
