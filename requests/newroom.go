package requests

import (
	"fmt"
	"log"

	obj "github.com/Brawdunoir/goplay-server/objects"
	"github.com/gorilla/websocket"
)

type NewRoomRequest struct {
	RoomName  string `json:"roomName"`
	OwnerName string `json:"ownerName"`
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
func (r NewRoomRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	// Retrieve owner info
	owner, err := users.User(r.OwnerName, remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("%w, cannot retrieve user info from database, has he logged in first ?", err)
	}

	roomID := rooms.AddRoom(r.RoomName, owner)

	log.Println(remoteAddr, "NewRoomRequest success")

	return roomID, nil
}
