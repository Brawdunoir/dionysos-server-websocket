package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// NewRoomRequest creates a new room within the server and
// send a NewRoomResponse has a confirmation for the creation.
type NewRoomRequest struct {
	RoomName  string `json:"roomName"`
	IsPrivate bool   `json:"isPrivate"`
}

func (r NewRoomRequest) Check() error {
	var err error

	if r.RoomName == "" {
		err = fmt.Errorf("%w; roomName is empty", err)
	}

	return err
}

// Handles a new room demand from a client.
func (r NewRoomRequest) Handle(publicAddr, uuid string, _ *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	// Retrieve owner info
	user, err := users.User(uuid, publicAddr, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	room := rooms.AddRoom(r.RoomName, user, r.IsPrivate, logger)

	logger.Infow("new room request", "owner", user.ID, "ownername", user.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.NewRoomResponse{RoomID: room.ID, RoomName: r.RoomName}, logger)
	return
}

func (r NewRoomRequest) Code() CodeType {
	return NEW_ROOM
}

func createNewRoomRequest(payload json.RawMessage) (r NewRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
