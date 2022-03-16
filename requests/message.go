package requests

import (
	"encoding/json"
	"fmt"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// NewMessageRequest send a message to all peers in the room.
// No message is store on the server side, the message must be stored
// on client side.
type NewMessageRequest struct {
	Salt    string `json:"salt"`
	RoomID  string `json:"roomId"`
	Content string `json:"content"`
}

func (r NewMessageRequest) Check() error {
	var err error

	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}
	if r.RoomID == "" {
		err = fmt.Errorf("%w; roomId is empty", err)
	}
	if r.Content == "" {
		err = fmt.Errorf("%w; content is empty", err)
	}

	return err
}

// Handles a new message from a client by forwarding it to all peers.
func (r NewMessageRequest) Handle(publicAddr string, _ *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, sender *obj.User) {
	// Fetch sender and room info
	sender, room, err := getUserAndRoom(r.Salt, publicAddr, r.RoomID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Check wethever the sender is in the room
	if !room.IsPeerPresent(sender, logger) {
		response = res.NewErrorResponse(constants.ERR_NO_PERM, logger)
		return
	}

	// Send the message to all peers in the room
	mes := obj.NewMessage(sender.ID, sender.Name, r.Content)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("new message request", "user", sender.ID, "username", sender.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r NewMessageRequest) Code() CodeType {
	return NEW_MESSAGE
}

func createNewMessageRequest(payload json.RawMessage) (r NewMessageRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
