package requests

import (
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

type NewMessageRequest struct {
	Username string `json:"username"`
	RoomID   string `json:"roomId"`
	Content  string `json:"content"`
}

func (r NewMessageRequest) Check() error {
	var err error

	if r.Username == "" {
		err = fmt.Errorf("%w; username is empty", err)
	}
	if r.Content == "" {
		err = fmt.Errorf("%w; content is empty", err)
	}

	return err
}

// Handles a new message from a client by forwarding it to all peers.
func (r NewMessageRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	// Fetch sender and room info
	sender, err := users.User(r.Username, remoteAddr)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	// Check wethever the sender is in the room
	if !room.IsPeerPresent(sender) {
		return res.NewErrorResponse(obj.NO_PERM)
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

	log.Println(remoteAddr, "NewMessageRequest success")

	return res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())})
}

func (r NewMessageRequest) Code() string {
	return NEW_MESSAGE
}
