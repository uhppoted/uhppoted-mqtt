package device

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"github.com/uhppoted/uhppoted-mqtt/common"
)

func (d *Device) GetCards(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.GetCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetCards(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve cards from %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) DeleteCards(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	rq := uhppoted.DeleteCardsRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.DeleteCards(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not delete cards on %d", *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) GetCard(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID   *uhppoted.DeviceID `json:"device-id"`
		CardNumber *uint32            `json:"card-number"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.CardNumber == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing card number", nil), fmt.Errorf("Invalid/missing card number")
	}

	rq := uhppoted.GetCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, err := impl.GetCard(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not retrieve card %v from %d", *body.CardNumber, *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) PutCard(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Card     *types.Card        `json:"card"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.Card == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing card", nil), fmt.Errorf("Invalid/missing card")
	}

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card:     *body.Card,
	}

	response, err := impl.PutCard(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not store card %v to %d", body.Card.CardNumber, *body.DeviceID), err), err
	}

	return response, nil
}

func (d *Device) DeleteCard(impl *uhppoted.UHPPOTED, request []byte) (interface{}, error) {
	body := struct {
		DeviceID   *uhppoted.DeviceID `json:"device-id"`
		CardNumber *uint32            `json:"card-number"`
	}{}

	if response, err := unmarshal(request, &body); err != nil {
		return response, err
	}

	if body.DeviceID == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing device ID", nil), fmt.Errorf("Invalid/missing device ID")
	}

	if body.CardNumber == nil {
		return common.MakeError(uhppoted.StatusBadRequest, "Invalid/missing card number", nil), fmt.Errorf("Invalid/missing card number")
	}

	rq := uhppoted.DeleteCardRequest{
		DeviceID:   *body.DeviceID,
		CardNumber: *body.CardNumber,
	}

	response, err := impl.DeleteCard(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not store card %v to %d", *body.CardNumber, *body.DeviceID), err), err
	}

	return response, nil
}
