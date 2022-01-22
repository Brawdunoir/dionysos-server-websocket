package requests

import (
	"errors"
	"fmt"
	"log"

	obj "github.com/Brawdunoir/goplay-server/objects"
	res "github.com/Brawdunoir/goplay-server/responses"
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
func (r JoinRoomAnswerRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	var v interface{}
	var err error

	if r.Accepted {
		v, err = handleAccept(r, remoteAddr, conn, users, rooms)
	} else {
		v, err = handleDeny(r, remoteAddr, conn, users, rooms)
	}

	log.Println(remoteAddr, "JoinRoomAnswerRequest success")

	return v, err
}

func handleAccept(r JoinRoomAnswerRequest, remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	// Fetch requester, owner and room info
	requester, err := users.UserByID(r.RequesterID)
	if err != nil {
		return nil, err
	}

	owner, err := users.User(r.OwnerName, remoteAddr)
	if err != nil {
		return nil, err
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		return nil, err
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != owner.ID {
		return nil, errors.New("you do not have this permission")
	}

	peers, err := rooms.Peers(r.RoomID)
	if err != nil {
		return nil, err
	}
	// Add the newcoming to the list of the peer before sending all the messages to existing peers
	rooms.AddPeer(r.RoomID, requester)

	// Send requester info to all actual peers
	for _, peer := range peers {
		peer.ConnMutex.Lock()
		peer.Conn.WriteJSON(res.Response{Status: res.NEW_PEER, RequestCode: JOIN_ROOM_ANSWER, Payload: res.Payload{Info: requester}})
		peer.ConnMutex.Unlock()
	}

	// Send actual peers info to the newcoming
	requester.ConnMutex.Lock()
	requester.Conn.WriteJSON(res.Response{Status: res.SUCCESS, RequestCode: JOIN_ROOM, Payload: res.Payload{Info: peers}})
	requester.ConnMutex.Unlock()

	return nil, nil
}

func handleDeny(r JoinRoomAnswerRequest, remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	requester, err := users.UserByID(r.RequesterID)
	if err != nil {
		return nil, err
	}

	requester.ConnMutex.Lock()
	requester.Conn.WriteJSON(res.Response{Status: res.REQUEST_DENIED, RequestCode: JOIN_ROOM})
	requester.ConnMutex.Unlock()

	return nil, nil
}
