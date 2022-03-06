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

// NewConnectionRequest is first request to the server.
// It registers the user within the server.
// The salt acts like a secret shared between the client and the server.
// It will identify the client along with its IP. It must be static during
// a session.
type NewConnectionRequest struct {
	Username string `json:"username"`
	Salt     string `json:"salt"`
}

func (r NewConnectionRequest) Check() error {
	var err error

	if r.Username == "" {
		err = fmt.Errorf("%w; username is empty", err)
	}
	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}

	return err
}

// Handles a new connection from a client.
func (r NewConnectionRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	userID := users.AddUser(r.Username, publicAddr, r.Salt, conn)

	log.Println(proxyAddr, "NewConnectionRequest success")

	return res.NewResponse(res.ConnectionResponse{Username: r.Username, UserID: userID})
}

func (r NewConnectionRequest) Code() CodeType {
	return NEW_CONNECTION
}

func createNewConnectionRequest(payload json.RawMessage) (r NewConnectionRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
