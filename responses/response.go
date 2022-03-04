package responses

import (
	"encoding/json"
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
		err := "not a valid IResponse payload"
		log.Println("Response does not implement IResponse interface", err)
		return NewErrorResponse(err)
	}
	_, err := r.Marshal()
	if err != nil {
		log.Println("cannot marshal IResponse to JSON")
		return NewErrorResponse(err.Error())
	}

	return createResponse(r)
}

// NewErrorResponse return a well formatted response with a error payload.
func NewErrorResponse(err string) Response {
	return NewResponse(ErrorResponse{Error: err})
}

func createResponse(r IResponse) Response {
	payload, err := r.Marshal()
	if err != nil {
		log.Println("error in createResponse", err)
	}
	return Response{Code: r.Code(), Payload: payload}
}
