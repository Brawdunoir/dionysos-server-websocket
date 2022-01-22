package requests

import (
	obj "github.com/Brawdunoir/goplay-server/objects"
	"github.com/gorilla/websocket"
)

// Requests codes
const (
	// Register a new user to the server and return the user's ID in the payload.
	NEW_CONNECTION = "NCO"
	// Register a new room to the server and return the room's ID in the payload.
	NEW_ROOM = "NRO"
	// Ask to join a room, return nothing. The answer is sent after owner decision.
	JOIN_ROOM = "JRO"
	// Follow a JOIN_ROOM request. Grant or deny user access to the room, return nothing.
	JOIN_ROOM_ANSWER = "AUT"
)

type IRequest interface {
	Check() error
	Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error)
}
