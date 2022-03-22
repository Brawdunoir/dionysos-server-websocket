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
func (r NewConnectionRequest) Handle(publicAddr, uuid string, conn *websocket.Conn, users *obj.Users, _ *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	user = users.AddUser(r.Username, publicAddr, uuid, conn, logger)

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
