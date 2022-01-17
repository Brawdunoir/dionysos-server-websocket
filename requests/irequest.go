package requests

import (
	obj "github.com/Brawdunoir/goplay-server/objects"
	"github.com/gorilla/websocket"
)

const (
	NEWCONNECTION    = "NCO"
	NEWROOM          = "NRO"
	JOINROOM         = "JRO"
	ACCEPTUSERTOROOM = "AUT"
	DENYUSERTOROOM   = "DUT"
)

type IRequest interface {
	Check() error
	Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error)
}
