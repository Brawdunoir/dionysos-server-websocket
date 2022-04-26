package requests

import (
	"encoding/json"
	"math"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// LoadFileRequest sets fileMetadata for the room and send file metadata to other peers
type LoadFileRequest struct {
	Name string `json:"name" validate:"min=3,max=254"` // Name of file
	Size uint64 `json:"size" validate:"min=1"`         // Size of file, in bytes

}

// Handle a file metadata upload. Happens just before the room's owner starts to upload file
func (r LoadFileRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
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

	// Set the fileMetadata and forward metadata to other peers
	// Chunks is basically the number of chunks that the file is composed from
	// (i.e. file total size divided by a single chunk size)
	chunks := uint16(math.Ceil(float64(r.Size) / float64(constants.CHUNK_SIZE)))
	room.FileMetadata = obj.FileMetadata{Name: r.Name, Size: r.Size, ChunkSize: constants.CHUNK_SIZE, Chunks: chunks}

	mes := res.NewResponse(res.LoadFileResponse(room.FileMetadata), logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("load file request success", "fileName", r.Name, "fileSize", r.Size, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r LoadFileRequest) Code() CodeType {
	return LOAD_FILE
}

func createLoadFileRequest(payload json.RawMessage) (r LoadFileRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
