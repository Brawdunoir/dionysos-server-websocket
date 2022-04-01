package requests

import (
	"encoding/json"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// ChangeUsernameRequest allows a user to change its username.
// The changement will trigger a RoomUserList update.
type ChangeUsernameRequest struct {
	NewUsername string `json:"newUsername" validate:"min=3,max=20"`
}

// Handles a username changement request from a client.
func (r ChangeUsernameRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch user an rename
	err := users.ChangeUsername(client.ID, r.NewUsername, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// If connected to a room, notify peers of the changement
	if client.RoomID != "" {
		room, err := rooms.Room(client.RoomID, logger)
		if err != nil {
			response = res.NewErrorResponse(err.Error(), logger)
			return
		}

		notifyPeers(rooms, room, logger)
	}

	logger.Infow("change username request", "user", client.ID, "username", client.Name)

	response = res.NewResponse(res.ChangeUsernameResponse{Username: r.NewUsername}, logger)
	return
}

func (r ChangeUsernameRequest) Code() CodeType {
	return CHANGE_USERNAME
}

func createChangeUsernameRequest(payload json.RawMessage) (r ChangeUsernameRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
