package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// QuitRoomRequest is for removing a user from a room.
type QuitRoomRequest struct {
	Username string `json:"username"`
}

func (r QuitRoomRequest) Check() error {
	var err error

	if r.Username == "" {
		err = fmt.Errorf("%w; username is empty", err)
	}

	return err
}

// Handles a quit request from a client.
// It removes the user from the room and it destroys the room if the room is empty.
// If the room is not empty it notify the remaining peers with an updated list of peers.
func (r QuitRoomRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch client and room info
	err := removeUserFromRoom(client, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

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
