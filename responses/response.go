package responses

import (
	"encoding/json"
)

// Status is a const defined in iresponse.go
// RequestCode is the request's code that triggered this response
// Payload carry more information, it can be empty
type Response struct {
	Status      string  `json:"status"`
	RequestCode string  `json:"requestcode"`
	Payload     Payload `json:"payload"`
}

// Payload carry more information, it can be empty
type Payload struct {
	Info json.RawMessage `json:"info,omitempty"`
	Err  string          `json:"error,omitempty"`
}

// CreateResponse return a SUCCESS/ERROR response and keep error from signature
func CreateResponse(info interface{}, requestCode string, err error) (Response, error) {
	var serr string
	if err == nil {
		serr = ""
	} else {
		serr = err.Error()
	}

	res, nerr := NewResponse(SUCCESS, requestCode, serr, info)
	if nerr != nil {
		return res, nerr
	}
	return res, err
}

// NewResponse return a simple response without keeping error from signature
func NewResponse(status, reqCode, err string, info interface{}) (Response, error) {
	var payload Payload

	if info != nil {
		info, err := json.Marshal(info)
		if err == nil {
			payload.Info = info
		} else {
			return Response{}, err
		}
	}
	if err != "" {
		payload.Err = err
	}
	return Response{Status: status, RequestCode: reqCode, Payload: payload}, nil
}
