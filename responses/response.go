package responses

import (
	"encoding/json"
	"errors"
	"log"
)

// Code is a const defined in iresponse.go
// Payload carry a specific response in this package
type Response struct {
	Code    CodeType        `json:"code"`
	Payload json.RawMessage `json:"payload"`
}

// NewResponse return a well formatted response.
func NewResponse(r IResponse) Response {
	if r.Code() == "" {
		err := errors.New("not a valid IResponse payload")
		log.Println("Response does not implement IResponse interface", err)
		return createResponse(ErrorResponse{Error: err})
	}
	_, err := r.MarshalJSON()
	if err != nil {
		log.Println("cannot marshal IResponse to JSON")
		return createResponse(ErrorResponse{Error: err})
	}

	return createResponse(r)
}

// NewErrorResponse return a well formatted response with a error payload.
func NewErrorResponse(err error) Response {
	return NewResponse(ErrorResponse{Error: err})
}

func createResponse(r IResponse) Response {
	payload, err := r.MarshalJSON()
	if err != nil {
		log.Println("error in createResponse", err)
	}
	return Response{Code: r.Code(), Payload: payload}
}
