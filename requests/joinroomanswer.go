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
func (r JoinRoomAnswerRequest) Handle(publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {

	if r.Accepted {
		response = handleAccept(r, publicAddr, proxyAddr, conn, users, rooms, logger)
	} else {
		response = handleDeny(r, publicAddr, proxyAddr, conn, users, rooms, logger)
	}

	return
}

func handleAccept(r JoinRoomAnswerRequest, publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	// Fetch requester, owner and room info
	requester, err := users.UserByID(r.RequesterID, logger)
	if err != nil {
		return res.NewErrorResponse("you are not connected", logger)
	}

	owner, err := users.User(r.OwnerSalt, publicAddr, logger)
	if err != nil {
		return res.NewErrorResponse("room's owner is disconnected", logger)
	}

	room, err := rooms.Room(r.RoomID, logger)
	if err != nil {
		return res.NewErrorResponse("the room does not exist or has been deleted", logger)
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != owner.ID {
		return res.NewErrorResponse("you do not have this permission, you are not the room's owner", logger)
	}

	// Add new peer to the list and notify all members
	err = addPeerAndNotify(requester, rooms, room, logger)
	if err != nil {
		return res.NewErrorResponse(err.Error(), logger)
	}

	logger.Infow("join room request", "user", requester.ID, "username", requester.Name, "owner", owner.ID, "ownername", owner.Name, "room", room.ID, "roomname", room.Name)

	return res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER}, logger)
}

func handleDeny(r JoinRoomAnswerRequest, publicAddr, proxyAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	requester, err := users.UserByID(r.RequesterID, logger)
	if err != nil {
		return res.NewErrorResponse("you are not connected", logger)
	}

	requesterResponse := res.NewResponse(res.DeniedResponse{RequestCode: JOIN_ROOM}, logger)
	requester.ConnMutex.Lock()
	requester.Conn.WriteJSON(requesterResponse)
	requester.ConnMutex.Unlock()

	logger.Infow("join room request", "user", requester.ID, "username", requester.Name, "room", r.RoomID)

	return res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER}, logger)
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
	for _, peer := range peers {
		peer.ConnMutex.Lock()
		peer.Conn.WriteJSON(res.NewResponse(res.NewPeersResponse{Peers: peers, OwnerID: room.OwnerID}, logger))
		peer.ConnMutex.Unlock()
		logger.Debugw("peer has been notified of peer change", "peer", peer.ID, "peername", peer.Name, "room", room.ID, "roomname", room.Name)
	}

	logger.Infow("notify peers", "room", room.ID, "roomname", room.Name)

	return nil
}
