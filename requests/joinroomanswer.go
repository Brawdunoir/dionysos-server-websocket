package requests

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// JoinRoomAnswerRequest indicates if a user (Requester) is
// accepted or not in the room
type JoinRoomAnswerRequest struct {
	RequesterID string `json:"requesterId" validate:"len=40"`
	Accepted    bool   `json:"accepted"`
}

// Grant or refuse access to room.
// In the first case, add the requester to the room and signal
// to every other peer in the room the newcoming, in addition to
// send the complete list of peer to the requester.
// In the second case, signal to the requester that his request had been denied
func (r JoinRoomAnswerRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {

	if r.Accepted {
		return handleAccept(r, client, users, rooms, logger)
	} else {
		return handleDeny(r, client, users, rooms, logger)
	}
}

func handleAccept(r JoinRoomAnswerRequest, client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch requester and room info
	requester, room, err := getUserByIdAndRoom(r.RequesterID, client.RoomID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != client.ID {
		response = res.NewErrorResponse(constants.ERR_NO_PERM, logger)
		return
	}

	// Add new peer to the list and notify all members
	err = addPeerAndNotify(requester, rooms, room, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	logger.Infow("join room request success", "user", requester.ID, "username", requester.Name, "owner", client.ID, "ownername", client.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func handleDeny(r JoinRoomAnswerRequest, client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	requester, err := users.UserByID(r.RequesterID, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	requesterResponse := res.NewResponse(res.DeniedResponse{RequestCode: JOIN_ROOM}, logger)
	requester.SendJSON(requesterResponse, logger)

	logger.Infow("join room request success", "user", requester.ID, "username", requester.Name, "owner", client.ID, "ownername", client.Name, "room", client.RoomID)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r JoinRoomAnswerRequest) Code() CodeType {
	return JOIN_ROOM_ANSWER
}

func createJoinRoomAnswerRequest(payload json.RawMessage) (r JoinRoomAnswerRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
