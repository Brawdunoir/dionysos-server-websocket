package responses

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
	Info interface{} `json:"info"`
	Err  string      `json:"error"`
}

// CreateResponse returns a response
func CreateResponse(pl interface{}, requestCode string, err error) (Response, error) {
	var status string
	var payload Payload

	if pl != nil || err != nil {
		payload = Payload{Info: pl, Err: err.Error()}
	}

	if err == nil {
		status = SUCCESS
	} else {
		status = ERROR
	}

	return Response{Status: status, RequestCode: requestCode, Payload: payload}, err
}
