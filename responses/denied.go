package responses

import "encoding/json"

type DeniedResponse struct {
	RequestCode CodeType `json:"requestCode"`
}

func (r DeniedResponse) Code() CodeType {
	return DENIED
}

func (r DeniedResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
