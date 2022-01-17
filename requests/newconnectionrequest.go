package requests

import (
	"fmt"
	"log"

	obj "github.com/Brawdunoir/goplay-server/objects"
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
func (r NewConnectionRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	userID := users.AddUser(r.Username, remoteAddr, conn)

	log.Println(remoteAddr, "NewConnectionRequest success")

	return userID, nil
}
