package go_ernie

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func (e *APIError) Error() string {
	if e.ErrorCode > 0 {
		return fmt.Sprintf("error, code: %d, message: %s", e.ErrorCode, e.ErrorMsg)
	}
	return e.ErrorMsg
}

func (e *APIError) UnmarshalJSON(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(data, &rawMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(rawMap["error_msg"], &e.ErrorMsg)
	if err != nil {
		return
	}
	return json.Unmarshal(rawMap["error_code"], &e.ErrorCode)
}

type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}
