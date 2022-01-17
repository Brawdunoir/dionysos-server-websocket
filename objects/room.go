package objects

import (
	"errors"
	"log"
	"sync"
)

// room represents data about a room for peers.
type room struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	OwnerID string     `json:"ownerid"`
	Peers   []*user    `json:"peers"`
	mu      sync.Mutex `json:"-"`
}

func (r *room) String() string {
	return r.Name + " (" + r.ID + ")"
}

// AddPeer safely adds a peer to a room
func (r *room) addPeer(u *user) error {
	if ok := r.IsPeerPresent(u); ok {
		return errors.New("peer already exists in room")
	}

	r.mu.Lock()
	r.Peers = append(r.Peers, u)
	r.mu.Unlock()
	return nil
}

// RemovePeer safely removes a peer from a room
func (r *room) removePeer(u *user) {
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

func (r *room) IsPeerPresent(u *user) bool {
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
func newRoom(roomName string, owner *user) *room {
	return &room{ID: generateRoomID(roomName, owner.RemoteAddr), Name: roomName, Peers: []*user{owner}}
}
