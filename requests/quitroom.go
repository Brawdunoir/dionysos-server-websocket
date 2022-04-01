package requests

import (
	"encoding/json"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// QuitRoomRequest is for removing a user from a room.
type QuitRoomRequest struct {
}

// Handles a quit request from a client.
// It removes the user from the room and it destroys the room if the room is empty.
// If the room is not empty it notify the remaining peers with an updated list of peers.
func (r QuitRoomRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch client and room info
	logger.Infow("user", client.ID, "username", client.Name, "room", client.RoomID)

	err := removeUserFromRoom(client, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	logger.Infow("quit room request success", "owner", client.ID, "ownername", client.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r QuitRoomRequest) Code() CodeType {
	return QUIT_ROOM
}

func createQuitRoomRequest(payload json.RawMessage) (r QuitRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
