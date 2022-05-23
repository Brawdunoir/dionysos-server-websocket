package responses

import "encoding/json"

// UploadChunkResponse is the counterpart of UploadChunkRequest and is forwarded to peers in the room
type UploadChunkResponse struct {
	ChunkNumber uint16 `json:"chunkNumber" validate:"nonzero"`
	ChunkData   []byte `json:"chunkData"` // Chunk data is validated during request in the handler
}

func (r UploadChunkResponse) Code() CodeType {
	return UPLOAD_CHUNK
}

func (r UploadChunkResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
