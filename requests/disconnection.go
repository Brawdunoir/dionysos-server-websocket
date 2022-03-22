package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// DisconnectionRequest is last request to the server.
// It unregisters the user within the server.
// The salt acts like a secret shared between the client and the server.
// It will identify the client along with its IP. It must be static during
// a session.
type DisconnectionRequest struct {
	Salt string `json:"salt"`
}

func (r DisconnectionRequest) Check() error {
	var err error

	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}

	return err
}

// Handles a new connection from a client.
func (r DisconnectionRequest) Handle(publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	user, err := users.User(r.Salt, publicAddr, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	users.RemoveUser(user.ID, logger)
	if user.RoomID != "" {
		rooms.RemovePeer(user.RoomID, user, logger)
	}

	logger.Infow("disconnection request", "user", user.ID, "username", user.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r DisconnectionRequest) Code() CodeType {
	return DISCONNECTION
}

func createDisconnectionRequest(payload json.RawMessage) (r DisconnectionRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
