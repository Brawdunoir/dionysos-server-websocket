package requests

import (
	"encoding/json"
	"errors"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// JoinRoomAnswerRequest indicates if a user (Requester) is
// accepted or not in the room
type JoinRoomAnswerRequest struct {
	OwnerSalt   string `json:"ownerSalt"`
	RoomID      string `json:"roomId"`
	RequesterID string `json:"requesterId"`
	Accepted    bool   `json:"accepted"`
}

func (r JoinRoomAnswerRequest) Check() error {
	var err error

	if r.OwnerSalt == "" {
		err = fmt.Errorf("%w; ownerSalt is empty", err)
	}
	if r.RoomID == "" {
		err = fmt.Errorf("%w; roomId is empty", err)
	}
	if r.RequesterID == "" {
		err = fmt.Errorf("%w; requesterId is empty", err)
	}

	return err
}

// Grant or refuse access to room.
// In the first case, add the requester to the room and signal
// to every other peer in the room the newcoming, in addition to
// send the complete list of peer to the requester.
// In the second case, signal to the requester that his request had been denied
func (r JoinRoomAnswerRequest) Handle(publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {

	if r.Accepted {
		return handleAccept(r, publicAddr, conn, users, rooms, logger)
	} else {
		return handleDeny(r, publicAddr, conn, users, rooms, logger)
	}
}

func handleAccept(r JoinRoomAnswerRequest, publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, owner *obj.User) {
	// Fetch requester, owner and room info
	requester, err := users.UserByID(r.RequesterID, logger)
	if err != nil {
		response = res.NewErrorResponse("you are not connected", logger)
		return

	}

	owner, err = users.User(r.OwnerSalt, publicAddr, logger)
	if err != nil {
		response = res.NewErrorResponse("room's owner is disconnected", logger)
		return
	}

	room, err := rooms.Room(r.RoomID, logger)
	if err != nil {
		response = res.NewErrorResponse("the room does not exist or has been deleted", logger)
		return
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != owner.ID {
		response = res.NewErrorResponse("you do not have this permission, you are not the room's owner", logger)
		return
	}

	// Add new peer to the list and notify all members
	err = addPeerAndNotify(requester, rooms, room, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	logger.Infow("join room request", "user", requester.ID, "username", requester.Name, "owner", owner.ID, "ownername", owner.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER}, logger)
	return
}

func handleDeny(r JoinRoomAnswerRequest, publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, owner *obj.User) {
	owner, err := users.User(r.OwnerSalt, publicAddr, logger)
	if err != nil {
		response = res.NewErrorResponse("room's owner is disconnected", logger)
		return
	}

	requester, err := users.UserByID(r.RequesterID, logger)
	if err != nil {
		response = res.NewErrorResponse("you are not connected", logger)
		return
	}

	requesterResponse := res.NewResponse(res.DeniedResponse{RequestCode: JOIN_ROOM}, logger)
	requester.SendJSON(requesterResponse, logger)

	logger.Infow("join room request", "user", requester.ID, "username", requester.Name, "owner", owner.ID, "ownername", owner.Name, "room", r.RoomID)

	response = res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER}, logger)
	return
}

func (r JoinRoomAnswerRequest) Code() CodeType {
	return JOIN_ROOM_ANSWER
}

func createJoinRoomAnswerRequest(payload json.RawMessage) (r JoinRoomAnswerRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}

func addPeerAndNotify(requester *obj.User, rooms *obj.Rooms, room *obj.Room, logger *zap.SugaredLogger) error {
	// Add the newcoming to the list of the peer before notifying
	_, err := rooms.AddPeer(room.ID, requester, logger)
	if err != nil {
		return err
	}

	err = notifyPeers(rooms, room, logger)
	if err != nil {
		return err
	}

	return nil
}

// Notify peers that the room peer list has changed
func notifyPeers(rooms *obj.Rooms, room *obj.Room, logger *zap.SugaredLogger) error {
	peers, err := rooms.Peers(room.ID, logger)
	if err != nil {
		return errors.New("error when retrieving peers in room")
	}

	// Send updated peers list to all peers
	mes := res.NewResponse(res.NewPeersResponse{Peers: peers, OwnerID: room.OwnerID}, logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("notify peers", "room", room.ID, "roomname", room.Name)

	return nil
}
