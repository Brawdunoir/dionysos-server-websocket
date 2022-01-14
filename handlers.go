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

	// // ReadJSON from client
	// err := conn.ReadJSON(&newConnectionRequest)
	// if err != nil {
	// 	return errors.New("cannot read JSON from client")
	// }

	// Create a new user and add it to the list
	user := NewUser(newConnectionRequest.Username, remoteAddr, conn)

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
	var newRoomRequest NewRoomRequest

	// // ReadJSON from client
	// err := conn.ReadJSON(&newRoomRequest)
	// if err != nil {
	// 	return errors.New("cannot read JSON from client")
	// }

	// Retrieve owner info
	ownerID := GenerateUserID(remoteAddr, newRoomRequest.OwnerName)
	owner, ok := users[ownerID]
	if !ok {
		return errors.New("cannot retrieve user info from database, has he logged in first ?")
	}

	room := NewRoom(newRoomRequest.RoomName, owner)

	// Add the room to map
	if _, exists := rooms[room.ID]; exists {
		log.Println("room already exists", room)
	} else {
		rooms[room.ID] = room
		log.Println("room successfully added", room)
	}

	return nil
}

func HandleJoinRoom(remoteAddr string, conn *websocket.Conn) error {
	var joinRoomRequest JoinRoomRequest

	// ReadJSON from client
	err := conn.ReadJSON(&joinRoomRequest)
	if err != nil {
		return errors.New("cannot read JSON from client")
	}

	// user, ok := users[joinRoomRequest.Username]
	// if !ok {
	// 	return errors.New("the user does not exist in the database")
	// }

	room, ok := rooms[joinRoomRequest.RoomID]
	if !ok {
		return errors.New("the room id does not match any existing room")
	}
	owner, ok := users[room.OwnerID]
	if !ok {
		return errors.New("the owner of the room " + room.ID + " does not exist anymore")
	}

	// * envoyer message au propriétaire
	owner.ConnMutex.Lock()
	// owner.Conn.WriteMessage(owner)
	// * recevoir sa réponse dans sa goroutine
	// * débloquer cette goroutine avec le channel du demandeur Accepted
	// * poursuivre l’exécution

	return nil
}
