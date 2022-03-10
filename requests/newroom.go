package requests

import (
	"encoding/json"
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// NewRoomRequest creates a new room within the server and
// send a NewRoomResponse has a confirmation for the creation.
type NewRoomRequest struct {
	RoomName  string `json:"roomName"`
	Salt      string `json:"salt"`
	IsPrivate bool   `json:"isPrivate"`
}

func (r NewRoomRequest) Check() error {
	var err error

	if r.RoomName == "" {
		err = fmt.Errorf("%w; roomName is empty", err)
	}
	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}

	return err
}

// Handles a new room demand from a client.
func (r NewRoomRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	// Retrieve owner info
	owner, err := users.User(r.Salt, publicAddr, logger)
	if err != nil {
		return res.NewErrorResponse(fmt.Sprintf("%w, cannot retrieve user info from database, has he logged in first ?", err))
	}

	roomID := rooms.AddRoom(r.RoomName, owner, r.IsPrivate, logger)

	log.Println(proxyAddr, "NewRoomRequest success")

	return res.NewResponse(res.NewRoomResponse{RoomID: roomID, RoomName: r.RoomName})
}

func (r NewRoomRequest) Code() CodeType {
	return NEW_ROOM
}

func createNewRoomRequest(payload json.RawMessage) (r NewRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
