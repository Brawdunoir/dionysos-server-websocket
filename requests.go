package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	NEWCONNECTION    = "NCO"
	NEWROOM          = "NRO"
	JOINROOM         = "JRO"
	ACCEPTUSERTOROOM = "AUT"
)

type IRequest interface {
	check() error
	handle(remoteAddr string, conn *websocket.Conn) (interface{}, error)
}

type Request struct {
	Code    string          `json:"code"`
	Payload json.RawMessage `json:"payload"`
}

type NewConnectionRequest struct {
	Username string `json:"username"`
}

type NewRoomRequest struct {
	RoomName  string `json:"roomname"`
	OwnerName string `json:"ownername"`
}

type JoinRoomRequest struct {
	Username string `json:"username"`
	RoomID   string `json:"roomid"`
}

type AcceptUserToRoomRequest struct {
	User User `json:"user"`
}

func (r Request) check() error {
	var err error

	if r.Code == "" {
		err = fmt.Errorf("%w; code is empty", err)
	}
	if !json.Valid(r.Payload) {
		err = fmt.Errorf("%w; json not valid", err)
	}

	return err
}

func (req Request) handle(remoteAddr string, conn *websocket.Conn) (interface{}, error) {
	err := req.check()
	if err != nil {
		return nil, err
	}

	var request IRequest

	// Would be better to change r type and unmarshall/handle this at the end of switch
	switch req.Code {
	case NEWCONNECTION:
		var r NewConnectionRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	case NEWROOM:
		var r NewRoomRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	// case JOINROOM:
	// 	var r JoinRoomRequest
	default:
		return nil, errors.New("unknown code")
	}
	if err != nil {
		return nil, err
	}

	return request.handle(remoteAddr, conn)
}

func (r NewConnectionRequest) check() error {
	var err error

	if r.Username == "" {
		err = fmt.Errorf("%w; username is empty", err)
	}

	return err
}

// Handles a new connection from a client.
// A client needs to send its name through a JSON after
// sending message NEWCONNECTION (see requests.go)
func (r NewConnectionRequest) handle(remoteAddr string, conn *websocket.Conn) (interface{}, error) {
	err := r.check()
	if err != nil {
		return nil, err
	}

	user := NewUser(r.Username, remoteAddr, conn)

	if _, exists := users[user.ID]; exists {
		log.Println("user already exists", user)
	} else {
		users[user.ID] = user
		log.Println("user successfully added", user)
	}

	return nil, nil
}

func (r NewRoomRequest) check() error {
	var err error

	if r.RoomName == "" {
		err = fmt.Errorf("%w; roomname is empty", err)
	}
	if r.OwnerName == "" {
		err = fmt.Errorf("%w; ownername is empty", err)
	}

	return err
}

// Handles a new room demand from a client.
// A client needs to send roomname and username through
// a JSON after sending message NEWROOM (see requests.go)
func (r NewRoomRequest) handle(remoteAddr string, conn *websocket.Conn) (interface{}, error) {
	err := r.check()
	if err != nil {
		return nil, err
	}

	// Retrieve owner info
	ownerID := GenerateUserID(remoteAddr, r.OwnerName)
	owner, ok := users[ownerID]
	if !ok {
		return nil, errors.New("cannot retrieve user info from database, has he logged in first ?")
	}

	room := NewRoom(r.RoomName, owner)

	// Add the room to map
	if _, exists := rooms[room.ID]; exists {
		log.Println("room already exists", room)
	} else {
		rooms[room.ID] = room
		log.Println("room successfully added", room)
	}

	return nil, nil
}

// func (r JoinRoomRequest) check() error {
// 	var err error

// 	if r.Username == "" {
// 		err = fmt.Errorf("%w; username is empty", err)
// 	}
// 	if r.RoomID == "" {
// 		err = fmt.Errorf("%w; roomid is empty", err)
// 	}

// 	return err
// }
