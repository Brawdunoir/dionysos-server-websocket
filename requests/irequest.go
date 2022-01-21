package requests

import (
	obj "github.com/Brawdunoir/goplay-server/objects"
	"github.com/gorilla/websocket"
)

// Requests codes
const (
	// Register a new user to the server and returns the user's ID in the payload
	NEW_CONNECTION = "NCO"
	// Register a new room to the server and returns the room's ID in the payload
	NEW_ROOM            = "NRO"
	JOIN_ROOM           = "JRO"
	ACCEPT_USER_TO_ROOM = "AUT"
	DENY_USER_TO_ROOM   = "DUT"
)

type IRequest interface {
	Check() error
	Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error)
}
