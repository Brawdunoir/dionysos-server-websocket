package responses

import "encoding/json"

// ErrorResponse correspond to a failed request
type ErrorResponse struct {
	Error string `json:"error"`
}

func (r ErrorResponse) Code() CodeType {
	return ERROR
}

func (r ErrorResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
