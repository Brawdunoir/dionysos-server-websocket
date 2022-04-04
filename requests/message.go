package requests

import (
	"encoding/json"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// NewMessageRequest send a message to all peers in the room.
// No message is store on the server side, the message must be stored
// on client side.
type NewMessageRequest struct {
	Content string `json:"content" validate:"nonzero"`
}

// Handles a new message from a client by forwarding it to all peers.
func (r NewMessageRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch client room info
	room, err := rooms.Room(client.RoomID, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Send the message to all peers in the room
	mes := res.NewResponse(res.NewMessageResponse{SenderID: client.ID, SenderName: client.Name, Content: r.Content}, logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("new message request success", "user", client.ID, "username", client.Name, "room", room.ID, "roomname", room.Name)

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
