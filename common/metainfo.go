package common

import (
	"encoding/json"
)

// TODO Remove - interim implemenation pending injection in MQTTD
type MetaInfo struct {
	RequestID *string `json:"request-id,omitempty"`
	ClientID  *string `json:"client-id,omitempty"`
	ServerID  string  `json:"server-id,omitempty"`
	Method    string  `json:"method,omitempty"`
	Nonce     Nonce   `json:"nonce,omitempty"`
}

type Nonce func() uint64

func (n Nonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(n())
}

// END TODO
