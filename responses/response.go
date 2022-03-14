package responses

import (
	"encoding/json"

	"go.uber.org/zap"
)

// Code is a const defined in iresponse.go
// Payload carry a specific response in this package
type Response struct {
	Code    CodeType        `json:"code"`
	Payload json.RawMessage `json:"payload"`
}

// NewResponse return a well formatted response.
func NewResponse(r IResponse, logger *zap.SugaredLogger) Response {
	if r.Code() == "" {
		err := "not a valid IResponse payload"
		logger.Errorw("response does not implement IResponse interface", "reponse", r)
		return NewErrorResponse(err, logger)
	}

	return createResponse(r, logger)
}

// NewErrorResponse return a well formatted response with a error payload.
func NewErrorResponse(err string, logger *zap.SugaredLogger) Response {
	return NewResponse(ErrorResponse{Error: err}, logger)
}

func createResponse(r IResponse, logger *zap.SugaredLogger) Response {
	payload, err := r.Marshal()
	if err != nil {
		logger.Error("cannot marshal IResponse to JSON")
		return NewErrorResponse(err.Error(), logger)
	}
	return Response{Code: r.Code(), Payload: payload}
}
