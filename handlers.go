package main

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

// Handles a new connection from a client.
// A client needs to send its name through a JSON after
// sending message NEWCONNECTION (see requests.go)
func HandleConnection(remoteAddr string, conn *websocket.Conn) error {
	var newConnectionRequest NewConnectionRequest

	// ReadJSON from client
	err := conn.ReadJSON(&newConnectionRequest)
	if err != nil || newConnectionRequest.Username == "" {
		return errors.New("wrong input from client")
	}

	// Create a new user and add it to the list
	user := NewUser(newConnectionRequest.Username, remoteAddr)

	if _, exists := users[user.ID]; exists {
		log.Println("user already exists", user)
	} else {
		users[user.ID] = user
		log.Println("user successfully added", user)
	}

	return nil
}

// Handles a new room demand from a client.
// A client needs to send roomname and username through
// a JSON after sending message NEWROOM (see requests.go)
func HandleNewRoom(remoteAddr string, conn *websocket.Conn) error {
	var newRoom NewRoomRequest

	// ReadJSON from client
	err := conn.ReadJSON(&newRoom)
	if err != nil || newRoom.RoomName == "" || newRoom.OwnerName == "" {
		return errors.New("wrong input from client")
	}

	// Retrieve owner info
	ownerID := GenerateUserID(remoteAddr, newRoom.OwnerName)
	owner, ok := users[ownerID]
	if !ok {
		return errors.New("cannot retrieve user info from database, has he logged in first ?")
	}

	room := NewRoom(newRoom.RoomName, owner)

	// Add the room to map
	if _, exists := rooms[room.ID]; exists {
		log.Println("room already exists", room)
	} else {
		rooms[room.ID] = room
		log.Println("room successfully added", room)
	}

	return nil
}
