package requests

import (
	"encoding/json"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
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
func (r DisconnectionRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	users.RemoveUser(client.ID, logger)
	if client.RoomID != "" {
		rooms.RemovePeer(client.RoomID, client, logger)
	}

	logger.Infow("disconnection request", "user", client.ID, "username", client.Name)

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
