package requests

import (
	"encoding/json"
	"fmt"
	"log"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

// JoinRoomRequest to server to ask permission for joining a room.
// The server then forward this request to the room's owner.
// RequesterID is set by the server and send to the room's owner.
// This way, the owner answer Yes or No to the request
// and join the RequesterID to his answer.
type JoinRoomRequest res.JoinRoomPendingResponse

func (r JoinRoomRequest) Check() error {
	var err error

	if r.RequesterUsername == "" {
		err = fmt.Errorf("%w; requesterUsername is empty", err)
	}
	if r.RoomID == "" {
		err = fmt.Errorf("%w; roomId is empty", err)
	}
	// RequesterID can be empty since it is replaced by server

	return err
}

// Handles a join room demand from a client by contacting
// the room's owner if the room is private.
// Otherwise it just add the peer to the room.
func (r JoinRoomRequest) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	// Fetch client, room and room owner info
	requester, err := users.User(r.RequesterUsername, remoteAddr)
	if err != nil {
		log.Println("can not retrieve requester info", err)
		return res.NewErrorResponse("you are not connected")
	}

	room, err := rooms.Room(r.RoomID)
	if err != nil {
		log.Println(err)
		return res.NewErrorResponse("the room does not exist or has been deleted")
	}
	owner, err := users.UserByID(room.OwnerID)
	if err != nil {
		log.Println("can not retrieve owner info", err)
		return res.NewErrorResponse("room's owner is disconnected")
	}

	if room.IsPeerPresent(requester) {
		return res.NewErrorResponse("you seem to be already in room")
	}

	if room.IsPrivate { // Private room: send request to room's owner for confirmation
		ownerRequest := res.NewResponse(res.JoinRoomPendingResponse{RoomID: room.ID, RequesterUsername: requester.Name, RequesterID: requester.ID})
		owner.ConnMutex.Lock()
		owner.Conn.WriteJSON(ownerRequest)
		owner.ConnMutex.Unlock()
	} else { // Public room: Directly add the peer and notify everybody
		addPeerAndNotify(requester, rooms, r.RoomID)
	}

	log.Println(remoteAddr, "JoinRoomRequest success")

	return res.NewResponse(res.JoinRoomResponse{RoomName: room.Name, RoomID: room.ID, IsPrivate: room.IsPrivate})
}

func (r JoinRoomRequest) Code() CodeType {
	return JOIN_ROOM
}

func createJoinRoomRequest(payload json.RawMessage) (r JoinRoomRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
