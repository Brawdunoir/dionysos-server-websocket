package responses

import "encoding/json"

type SuccessResponse struct {
	RequestCode CodeType `json:"requestCode"`
}

func (r SuccessResponse) Code() CodeType {
	return SUCCESS
}

func (r SuccessResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
