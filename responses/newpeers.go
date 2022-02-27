package responses

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/objects"
)

type NewPeersResponse struct {
	Peers objects.PeersType `json:"peers"`
}

func (r NewPeersResponse) Code() CodeType {
	return NEW_PEER
}

func (r NewPeersResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
