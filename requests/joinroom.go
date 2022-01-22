package requests

import (
	"errors"
	"fmt"
	"log"

	obj "github.com/Brawdunoir/goplay-server/objects"
	res "github.com/Brawdunoir/goplay-server/responses"
	"github.com/gorilla/websocket"
)

// JoinRoomRequest to server to ask permission for joining a room.
// The server then forward this request to the room's owner.
// RequesterID is set by the server and send to the room's owner.
// This way, the owner answer Yes or No to the request
// and join the RequesterID to his answer.
type JoinRoomRequest struct {
	RequesterUsername string `json:"requesterUsername"`
	RoomID            string `json:"roomId"`
	RequesterID       string `json:"requesterId,omitempty"`
}

func (r JoinRoomRequest) Check() error {
	var err error

	if r.RequesterUsername == "" {
		err = fmt.Errorf("%w; requesterUsername is empty", err)
	}
	if r.RoomID == "" {
		err = fmt.Errorf("%w; roomId is empty", err)
	}
	// RequesterID can be empty since it is replaced by server

	return err
}

// Handles a join room demand from a client by contacting the room's owner
func (r JoinRoomRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	// Fetch client, room and room owner info
	requester, err := users.User(r.RequesterUsername, remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("%w, can not retrieve requester info", err)
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		return nil, err
	}
	owner, err := users.UserByID(room.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("%w, can not retrieve owner info", err)
	}

	if room.IsPeerPresent(requester) {
		return nil, errors.New("user is already in the room")
	}

	// Add RequesterID to the request before sending it to the room's owner
	r.RequesterID = requester.ID

	owner.ConnMutex.Lock()
	owner.Conn.WriteJSON(res.NewResponse(res.JOIN_ROOM_PENDING, JOIN_ROOM, "", r))
	owner.ConnMutex.Unlock()

	log.Println(remoteAddr, "JoinRoomRequest success")

	return nil, nil
}
