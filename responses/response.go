package responses

import "encoding/json"

// Status is a const defined in iresponse.go
// RequestCode is the request's code that triggered this response
// Payload carry more information, it can be empty
type Response struct {
	Status      string          `json:"status"`
	RequestCode string          `json:"requestcode"`
	Payload     json.RawMessage `json:"payload"`
}

type payload struct {
	Info interface{} `json:"info"`
	Err  string      `json:"error"`
}

// CreateResponse returns a response
func CreateResponse(pl interface{}, err error, requestCode string) (Response, error) {
	jsonpayload, erro := json.Marshal(payload{Info: pl, Err: err.Error()})
	if erro != nil {
		return Response{}, erro
	}
	if err == nil {
		return Response{Status: SUCCESS, RequestCode: requestCode, Payload: jsonpayload}, nil
	} else {
		return Response{Status: ERROR, RequestCode: requestCode, Payload: jsonpayload}, nil
	}
}
