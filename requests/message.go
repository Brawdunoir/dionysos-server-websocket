package requests

import (
	"encoding/json"
	"fmt"
	"log"

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
func (r NewMessageRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	// Fetch sender and room info
	sender, err := users.User(r.Salt, publicAddr)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	// Check wethever the sender is in the room
	if !room.IsPeerPresent(sender) {
		return res.NewErrorResponse(constants.ERR_NO_PERM)
	}

	// Gather all peers and send the new message to all peers
	peers, err := rooms.Peers(r.RoomID)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	for _, peer := range peers {
		peer.ConnMutex.Lock()
		peer.Conn.WriteJSON(obj.NewMessage(sender.ID, sender.Name, r.Content))
		peer.ConnMutex.Unlock()
	}

	log.Println(proxyAddr, "NewMessageRequest success")

	return res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())})
}

func (r NewMessageRequest) Code() CodeType {
	return NEW_MESSAGE
}

func createNewMessageRequest(payload json.RawMessage) (r NewMessageRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
