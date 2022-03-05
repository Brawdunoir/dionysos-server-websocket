package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

type JoinRoomAnswerRequest struct {
	OwnerName   string `json:"ownerName"`
	RoomID      string `json:"roomId"`
	RequesterID string `json:"requesterId"`
	Accepted    bool   `json:"accepted"`
}

func (r JoinRoomAnswerRequest) Check() error {
	var err error

	if r.OwnerName == "" {
		err = fmt.Errorf("%w; ownerName is empty", err)
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
func (r JoinRoomAnswerRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (response res.Response) {

	if r.Accepted {
		response = handleAccept(r, remoteAddr, conn, users, rooms)
	} else {
		response = handleDeny(r, remoteAddr, conn, users, rooms)
	}

	log.Println(remoteAddr, "JoinRoomAnswerRequest success")

	return
}

func handleAccept(r JoinRoomAnswerRequest, remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	// Fetch requester, owner and room info
	requester, err := users.UserByID(r.RequesterID)
	if err != nil {
		log.Println("can not retrieve requester info", err)
		return res.NewErrorResponse("you are not connected")
	}

	owner, err := users.User(r.OwnerName, remoteAddr)
	if err != nil {
		log.Println("can not retrieve owner info", err)
		return res.NewErrorResponse("room's owner is disconnected")
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		log.Println(err)
		return res.NewErrorResponse("the room does not exist or has been deleted")
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != owner.ID {
		return res.NewErrorResponse("you do not have this permission, you are not the room's owner")
	}

	// Add new peer to the list and notify all members
	err = addPeerAndNotify(requester, rooms, r.RoomID)
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	return res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER})
}

func handleDeny(r JoinRoomAnswerRequest, remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	requester, err := users.UserByID(r.RequesterID)
	if err != nil {
		log.Println(err)
		return res.NewErrorResponse("you are not connected")
	}

	// res, err := res.NewResponse(res.REQUEST_DENIED, JOIN_ROOM, "", nil)
	requesterResponse := res.NewResponse(res.DeniedResponse{RequestCode: JOIN_ROOM})
	requester.ConnMutex.Lock()
	requester.Conn.WriteJSON(requesterResponse)
	requester.ConnMutex.Unlock()

	return res.NewResponse(res.SuccessResponse{RequestCode: JOIN_ROOM_ANSWER})
}

func (r JoinRoomAnswerRequest) Code() CodeType {
	return JOIN_ROOM_ANSWER
}

func addPeerAndNotify(requester *obj.User, rooms *obj.Rooms, roomID string) error {
	// Add the newcoming to the list of the peer before notifying
	rooms.AddPeer(roomID, requester)

	peers, err := rooms.Peers(roomID)
	if err != nil {
		log.Println(err)
		return errors.New("error when retrieving peers in room")
	}

	// Send updated peers list to all peers
	for _, peer := range peers {
		peer.ConnMutex.Lock()
		peer.Conn.WriteJSON(res.NewResponse(res.NewPeersResponse{Peers: peers}))
		peer.ConnMutex.Unlock()
	}

	return nil
}

func createJoinRoomAnswerRequest(payload json.RawMessage) (r JoinRoomAnswerRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
