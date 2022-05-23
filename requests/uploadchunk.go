package requests

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// UploadChunkRequest uploads a piece of file (chunk) to the server so it can forward it to other peers
type UploadChunkRequest struct {
	ChunkNumber uint16 `json:"chunkNumber" validate:"nonzero"`
	ChunkData   []byte `json:"chunkData"` // Chunk data is validated during request in the handler
}

// Handle a chunk upload by forwarding it to the other peers in the room.
// Validates that the size of chunk data is exactly the same as the constant CHUNK_SIZE.
// Except for the last chunk, which can have a lower size.
func (r UploadChunkRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch room info
	room, err := rooms.Room(client.RoomID, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != client.ID {
		response = res.NewErrorResponse(constants.ERR_NO_PERM, logger)
		return
	}

	// Validates chunk data length
	err = validateChunkData(room.FileMetadata, r)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Send the chunk to all peers in the room
	mes := res.NewResponse(res.UploadChunkResponse(r), logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("upload chunk request success", "chunkNumber", r.ChunkNumber, "fileName", room.FileMetadata.Name, "chunkSize", room.FileMetadata.ChunkSize, "fileSize", room.FileMetadata.Size, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r UploadChunkRequest) Code() CodeType {
	return UPLOAD_CHUNK
}

func createUploadChunkRequest(payload json.RawMessage) (r UploadChunkRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
