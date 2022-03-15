package requests

import (
	"encoding/json"
	"fmt"

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
func (r NewConnectionRequest) Handle(publicAddr string, conn *websocket.Conn, users *obj.Users, _ *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	user = users.AddUser(r.Username, publicAddr, r.Salt, conn, logger)

	logger.Infow("connection request", "user", user.ID, "username", user.Name)

	response = res.NewResponse(res.ConnectionResponse{Username: r.Username, UserID: user.ID}, logger)
	return
}

func (r NewConnectionRequest) Code() CodeType {
	return NEW_CONNECTION
}

func createNewConnectionRequest(payload json.RawMessage) (r NewConnectionRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
