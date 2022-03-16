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
func (r JoinRoomRequest) Handle(publicAddr string, _ *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	// Fetch client, room and room owner info
	user, owner, room, err := getUserAndRoomAndRoomOwner(r.Salt, publicAddr, r.RoomID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	if room.IsPeerPresent(user, logger) {
		response = res.NewErrorResponse("you seem to be already in room", logger)
		return
	}

	if room.IsPrivate { // Private room: send request to room's owner for confirmation
		ownerRequest := res.NewResponse(res.JoinRoomPendingResponse{RoomID: room.ID, RequesterUsername: user.Name, RequesterID: user.ID}, logger)
		owner.SendJSON(ownerRequest, logger)
	} else { // Public room: Directly add the peer and notify everybody
		err = addPeerAndNotify(user, rooms, room, logger)
		if err != nil {
			response = res.NewErrorResponse(err.Error(), logger)
			return
		}
	}

	logger.Infow("join room request", "user", user.ID, "username", user.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.JoinRoomResponse{RoomName: room.Name, RoomID: room.ID, IsPrivate: room.IsPrivate}, logger)
	return
}

func (r JoinRoomRequest) Code() CodeType {
	return JOIN_ROOM
}

func createJoinRoomRequest(payload json.RawMessage) (r JoinRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
