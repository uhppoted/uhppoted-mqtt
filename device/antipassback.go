package device

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetAntiPassback(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	body := struct {
		Controller uint32 `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.Controller == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing device ID")
	}

	if antipassback, err := impl.GetAntiPassback(body.Controller); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Could not retrieve antipassback for %v", body.Controller), err), err
	} else {
		response := struct {
			Controller   uint32 `json:"device-id"`
			AntiPassback string `json:"anti-passback"`
		}{
			Controller:   body.Controller,
			AntiPassback: fmt.Sprintf("%v", antipassback),
		}

		return response, nil
	}
}

func (d *Device) SetAntiPassback(impl uhppoted.IUHPPOTED, request []byte) (any, error) {
	normalise := func(v string) string {
		return strings.ToLower(regexp.MustCompile(`[ (),]+`).ReplaceAllString(v, ""))
	}

	parse := func(v string) (types.AntiPassback, error) {
		switch normalise(v) {
		case "disabled":
			return types.Disabled, nil

		case "1:2;3:4":
			return types.Readers12_34, nil

		case "13:24":
			return types.Readers13_24, nil

		case "1:23":
			return types.Readers1_23, nil

		case "1:234":
			return types.Readers1_234, nil

		default:
			return types.Disabled, fmt.Errorf("invalid/missing anti-passback")
		}
	}

	body := struct {
		Controller   uint32 `json:"device-id"`
		AntiPassback string `json:"anti-passback"`
	}{}

	if v, err := unmarshal(request, &body); err != nil {
		return v, err
	}

	if controller := body.Controller; controller == 0 {
		return common.MakeError(StatusBadRequest, "Invalid/missing controller ID", nil), fmt.Errorf("invalid/missing controller ID")
	} else if antipassback, err := parse(body.AntiPassback); err != nil {
		return common.MakeError(StatusBadRequest, "Invalid/missing anti-passback", nil), err
	} else if ok, err := impl.SetAntiPassback(body.Controller, antipassback); err != nil {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Error setting anti-passback for %v", controller), err), err
	} else if !ok {
		return common.MakeError(StatusInternalServerError, fmt.Sprintf("Failed to update anti-passback for %v", controller), nil), fmt.Errorf("failed to set anti-passback for %v", controller)
	} else {
		response := struct {
			Controller   uint32 `json:"device-id"`
			AntiPassback string `json:"anti-passback"`
		}{
			Controller:   controller,
			AntiPassback: fmt.Sprintf("%v", antipassback),
		}

		return response, nil
	}
}
