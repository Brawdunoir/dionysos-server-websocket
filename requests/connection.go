package requests

import (
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

type NewConnectionRequest struct {
	Username string `json:"username"`
}

func (r NewConnectionRequest) Check() error {
	var err error

	if r.Username == "" {
		err = fmt.Errorf("%w; username is empty", err)
	}

	return err
}

// Handles a new connection from a client.
func (r NewConnectionRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	userID := users.AddUser(r.Username, remoteAddr, conn)

	log.Println(remoteAddr, "NewConnectionRequest success")

	return res.NewResponse(res.ConnectionResponse{Username: r.Username, UserID: userID})
}

func (r NewConnectionRequest) Code() CodeType {
	return NEW_CONNECTION
}
