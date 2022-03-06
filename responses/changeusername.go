package responses

import "encoding/json"

// ChangeUsernameResponse correspond to a successful connection/registration
// of a user.
type ChangeUsernameResponse struct {
	Username string `json:"username"`
}

func (r ChangeUsernameResponse) Code() CodeType {
	return CONNECTION_SUCCESS
}

func (r ChangeUsernameResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
