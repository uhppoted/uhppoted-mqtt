package common

type Error struct {
	Code    int    `json:"error-code"`
	Message string `json:"message,omitempty"`
	Debug   error  `json:"debug,omitempty"`
}
