package common

type Error struct {
	Error struct {
		Code    int    `json:"error-code"`
		Message string `json:"message,omitempty"`
		Debug   string `json:"debug,omitempty"`
	} `json:"error"`
}

func MakeError(code int, message string, debug error) *Error {
	var dbg string

	if debug != nil {
		dbg = debug.Error()
	}

	err := struct {
		Code    int    `json:"error-code"`
		Message string `json:"message,omitempty"`
		Debug   string `json:"debug,omitempty"`
	}{
		Code:    code,
		Message: message,
		Debug:   dbg,
	}

	return &Error{
		Error: err,
	}
}
