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
	type card struct {
		CardNumber uint32                `json:"card-number"`
		From       *types.Date           `json:"start-date"`
		To         *types.Date           `json:"end-date"`
		Doors      map[uint8]interface{} `json:"doors"`
	}
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Card     *card              `json:"card"`
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

	c := *body.Card

	rq := uhppoted.PutCardRequest{
		DeviceID: *body.DeviceID,
		Card: types.Card{
			CardNumber: c.CardNumber,
			From:       c.From,
			To:         c.To,
			Doors:      map[uint8]int{1: 0, 2: 0, 3: 0, 4: 0},
		},
	}

	for _, k := range []uint8{1, 2, 3, 4} {
		switch v := c.Doors[k].(type) {
		case bool:
			if v {
				rq.Card.Doors[k] = 1
			}

		case int:
			if v >= 2 && v < 254 {
				rq.Card.Doors[k] = v
			}

		case float64:
			if int(v) >= 2 && int(v) < 254 {
				rq.Card.Doors[k] = int(v)
			}
		}
	}

	response, err := impl.PutCard(rq)
	if err != nil {
		return common.MakeError(uhppoted.StatusInternalServerError, fmt.Sprintf("Could not store card %v to %d", body.Card.CardNumber, *body.DeviceID), err), err
	}

	doors := map[uint8]interface{}{1: false, 2: false, 3: false, 4: false}
	for _, k := range []uint8{1, 2, 3, 4} {
		v := response.Card.Doors[k]
		switch v {
		case 0:
			doors[k] = false
		case 1:
			doors[k] = true
		default:
			doors[k] = v
		}
	}

	return struct {
		DeviceID uhppoted.DeviceID `json:"device-id"`
		Card     card              `json:"card"`
	}{
		DeviceID: response.DeviceID,
		Card: card{
			CardNumber: response.Card.CardNumber,
			From:       response.Card.From,
			To:         response.Card.To,
			Doors:      doors,
		},
	}, nil
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
