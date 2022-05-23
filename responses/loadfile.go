package responses

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/objects"
)

// LoadFileResponse correspond to a successful LoadFileRequest.
// It sends file metadata to peers in the room.
type LoadFileResponse objects.FileMetadata

func (r LoadFileResponse) Code() CodeType {
	return LOAD_FILE
}

func (r LoadFileResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
