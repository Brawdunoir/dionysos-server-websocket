package responses

import "encoding/json"

// ErrorResponse correspond to a failed request
type ErrorResponse struct {
	Error error `json:"error"`
}

func (r ErrorResponse) Code() CodeType {
	return ERROR
}

func (r ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}
