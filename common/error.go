package common

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Debug   string `json:"debug,omitempty"`
}

func MakeError(code int, message string, debug error) *Error {
	var dbg string

	if debug != nil {
		dbg = debug.Error()
	}

	return &Error{
		Code:    code,
		Message: message,
		Debug:   dbg,
	}
}
