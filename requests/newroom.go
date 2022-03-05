package requests

import (
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

type NewRoomRequest struct {
	RoomName  string `json:"roomName"`
	OwnerName string `json:"ownerName"`
	IsPrivate bool   `json:"isPrivate"`
}

func (r NewRoomRequest) Check() error {
	var err error

	if r.RoomName == "" {
		err = fmt.Errorf("%w; roomName is empty", err)
	}
	if r.OwnerName == "" {
		err = fmt.Errorf("%w; ownerName is empty", err)
	}

	return err
}

// Handles a new room demand from a client.
func (r NewRoomRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	// Retrieve owner info
	owner, err := users.User(r.OwnerName, remoteAddr)
	if err != nil {
		return res.NewErrorResponse(fmt.Sprintf("%w, cannot retrieve user info from database, has he logged in first ?", err))
	}

	roomID := rooms.AddRoom(r.RoomName, owner, r.IsPrivate)

	log.Println(remoteAddr, "NewRoomRequest success")

	return res.NewResponse(res.NewRoomResponse{RoomID: roomID, RoomName: r.RoomName})
}

func (r NewRoomRequest) Code() CodeType {
	return NEW_ROOM
}
