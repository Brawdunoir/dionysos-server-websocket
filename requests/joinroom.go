package requests

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// JoinRoomRequest to server to ask permission for joining a room.
// The server then forward this request to the room's owner.
// This way, the owner answer Yes or No to the request
// and join the RequesterID to his answer.
type JoinRoomRequest struct {
	RoomID string `json:"roomId" validate:"len=40"`
}

// Handles a join room demand from a client by contacting
// the room's owner if the room is private.
// Otherwise it just add the peer to the room.
func (r JoinRoomRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch client, room and room owner info
	room, owner, err := getRoomAndRoomOwner(r.RoomID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	if room.IsPeerPresent(client, logger) {
		response = res.NewErrorResponse(constants.ERR_PEER_ALREADY_IN_ROOM, logger)
		return
	}

	if room.IsPrivate { // Private room: send request to room's owner for confirmation
		ownerRequest := res.NewResponse(res.JoinRoomPendingResponse{RoomID: room.ID, RequesterUsername: client.Name, RequesterID: client.ID}, logger)
		owner.SendJSON(ownerRequest, logger)
	} else { // Public room: Directly add the peer and notify everybody
		err = addPeerAndNotify(client, rooms, room, logger)
		if err != nil {
			response = res.NewErrorResponse(err.Error(), logger)
			return
		}
	}

	logger.Infow("join room request success", "user", client.ID, "username", client.Name, "room", room.ID, "roomname", room.Name)

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
