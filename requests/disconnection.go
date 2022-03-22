package requests

import (
	"encoding/json"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// DisconnectionRequest is last request to the server.
// It unregisters the user within the server.
type DisconnectionRequest struct {
}

func (r DisconnectionRequest) Check() error {
	var err error

	return err
}

// Handles a new connection from a client.
func (r DisconnectionRequest) Handle(publicAddr, uuid string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	user, err := users.User(uuid, publicAddr, logger)
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
