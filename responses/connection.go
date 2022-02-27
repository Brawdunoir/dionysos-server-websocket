package responses

import "encoding/json"

// ConnectionResponse correspond to a successful connection/registration
// of a user.
type ConnectionResponse struct {
	Username string `json:"username"`
	UserID   string `json:"userId"`
}

func (r ConnectionResponse) Code() CodeType {
	return CONNECTION_SUCCESS
}

func (r ConnectionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
