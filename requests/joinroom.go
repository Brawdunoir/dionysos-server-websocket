package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// JoinRoomRequest to server to ask permission for joining a room.
// The server then forward this request to the room's owner.
// This way, the owner answer Yes or No to the request
// and join the RequesterID to his answer.
type JoinRoomRequest struct {
	RoomID string `json:"roomId"`
	Salt   string `json:"salt"`
}

func (r JoinRoomRequest) Check() error {
	var err error

	if r.RoomID == "" {
		err = fmt.Errorf("%w; roomId is empty", err)
	}
	if r.Salt == "" {
		err = fmt.Errorf("%w; salt is empty", err)
	}
	// RequesterID can be empty since it is replaced by server

	return err
}

// Handles a join room demand from a client by contacting
// the room's owner if the room is private.
// Otherwise it just add the peer to the room.
func (r JoinRoomRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	// Fetch client, room and room owner info
	requester, err := users.User(r.Salt, publicAddr, logger)
	if err != nil {
		return res.NewErrorResponse("you are not connected", logger)
	}

	room, err := rooms.Room(r.RoomID, logger)
	if err != nil {
		return res.NewErrorResponse("the room does not exist or has been deleted", logger)
	}
	owner, err := users.UserByID(room.OwnerID, logger)
	if err != nil {
		return res.NewErrorResponse("room's owner is disconnected", logger)
	}

	if room.IsPeerPresent(requester, logger) {
		return res.NewErrorResponse("you seem to be already in room", logger)
	}

	if room.IsPrivate { // Private room: send request to room's owner for confirmation
		ownerRequest := res.NewResponse(res.JoinRoomPendingResponse{RoomID: room.ID, RequesterUsername: requester.Name, RequesterID: requester.ID}, logger)
		owner.ConnMutex.Lock()
		owner.Conn.WriteJSON(ownerRequest)
		owner.ConnMutex.Unlock()
	} else { // Public room: Directly add the peer and notify everybody
		addPeerAndNotify(requester, rooms, room, logger)
	}

	logger.Infow("join room request", "user", requester.ID, "username", requester.Name, "room", room.ID, "roomname", room.Name)

	return res.NewResponse(res.JoinRoomResponse{RoomName: room.Name, RoomID: room.ID, IsPrivate: room.IsPrivate}, logger)
}

func (r JoinRoomRequest) Code() CodeType {
	return JOIN_ROOM
}

func createJoinRoomRequest(payload json.RawMessage) (r JoinRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
