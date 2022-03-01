package responses

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/objects"
)

type NewMessageResponse objects.Message

func (r NewMessageResponse) Code() CodeType {
	return NEW_MESSAGE
}

func (r NewMessageResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
