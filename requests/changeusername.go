package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ChangeUsernameRequest allows a user to change its username.
// The changement will trigger a RoomUserList update.
type ChangeUsernameRequest struct {
	NewUsername string `json:"newUsername"`
	Salt        string `json:"salt"`
}

func (r ChangeUsernameRequest) Check() error {
	var err error

	if r.NewUsername == "" {
		err = fmt.Errorf("%w; newUsername is empty", err)
	}
	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}

	return err
}

// Handles a username changement request from a client.
func (r ChangeUsernameRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	// Fetch user an rename
	user, err := users.User(r.Salt, publicAddr, logger)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	err = users.ChangeUsername(user.ID, r.NewUsername, logger)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	// If connected to a room, notify peers of the changement
	if user.RoomID != "" {
		room, err := rooms.Room(user.RoomID, logger)
		if err != nil {
			return res.NewErrorResponse(err.Error())
		}

		notifyPeers(rooms, room, logger)
	}

	logger.Infow("change username request", "user", user.ID, "username", user.Name)

	return res.NewResponse(res.ChangeUsernameResponse{Username: r.NewUsername})
}

func (r ChangeUsernameRequest) Code() CodeType {
	return CHANGE_USERNAME
}

func createChangeUsernameRequest(payload json.RawMessage) (r ChangeUsernameRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
