package main

import (
	"crypto/sha1"
	"fmt"
	"log"
)

// room represents data about a room for peers.
type Room struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"ownerid"`
	Peers   []User `json:"peers"`
}

func (r Room) String() string {
	return r.Name + " (" + r.ID + ")"
}

func (r *Room) AddPeer(u User) {
	for _, p := range r.Peers {
		if p.ID == u.ID {
			return
		}
	}
	r.Peers = append(r.Peers, u)
}

// Remove a User u from Room r
func (r *Room) RemovePeer(u User) {
	for i, p := range r.Peers {
		if p.ID == u.ID {
			r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
			log.Println("user", u, "removed from the room", r)
			return
		}
	}
	log.Println("can't find", u, "in room", r)
}

// Create a hash for a room
// Return an error if name or Owner is empty
func (r *Room) CreateID() error {
	r.ID = fmt.Sprint(sha1.Sum([]byte(r.Name)))

	return nil
}

// Generate a room ID based on a roomname and an ownername
func GenerateRoomID(roomName, ownerName string) string {
	return generateStringHash(roomName + ownerName)
}

// Creates a new room
func NewRoom(roomName string, owner User) Room {
	return Room{ID: GenerateRoomID(roomName, owner.Name), Name: roomName, Peers: []User{owner}}
}
