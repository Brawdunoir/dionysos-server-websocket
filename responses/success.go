package responses

// SuccessReponse is a specific response, it only echo its RequestCode.
type SuccessResponse struct {
	RequestCode CodeType `json:"requestCode"`
}

func (r SuccessResponse) Code() CodeType {
	return r.RequestCode
}

func (r SuccessResponse) Marshal() ([]byte, error) {
	return nil, nil
}
