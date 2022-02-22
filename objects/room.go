package objects

import (
	"errors"
	"log"
	"sync"
)

type PeersType []*User

// room represents data about a room for peers.
type Room struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	OwnerID string     `json:"ownerid"`
	Peers   PeersType  `json:"peers"`
	mu      sync.Mutex `json:"-"`
}

func (r *Room) String() string {
	return r.Name + " (" + r.ID + ")"
}

// AddPeer safely adds a peer to a room
func (r *Room) AddPeer(u *User) error {
	if ok := r.IsPeerPresent(u); ok {
		return errors.New("peer already exists in room")
	}

	r.mu.Lock()
	r.Peers = append(r.Peers, u)
	r.mu.Unlock()
	return nil
}

// RemovePeer safely removes a peer from a room
func (r *Room) RemovePeer(u *User) {
	for i, p := range r.Peers {
		if p.ID == u.ID {
			r.mu.Lock()
			r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
			r.mu.Unlock()
			log.Println("user", u, "removed from the room", r)
			return
		}
	}
	log.Println("can't find", u, "in room", r)
}

func (r *Room) IsPeerPresent(u *User) bool {
	for _, p := range r.Peers {
		if p.ID == u.ID {
			return true
		}
	}
	return false
}

// Generate a room ID based on a roomname and an ownerRemoteAddr
func generateRoomID(roomName, ownerRemoteAddr string) string {
	return generateStringHash(roomName + ownerRemoteAddr)
}

// Creates a new room
func NewRoom(roomName string, owner *User) *Room {
	return &Room{ID: generateRoomID(roomName, owner.RemoteAddr), Name: roomName, OwnerID: owner.ID, Peers: PeersType{owner}}
}
