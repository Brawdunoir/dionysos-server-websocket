package objects

import (
	"errors"
	"log"
	"sync"

	"github.com/Brawdunoir/dionysos-server/constants"
)

type Rooms struct {
	saloons map[string]*Room
	mu      sync.RWMutex
}

// Creates a new set of Rooms
func NewRooms() *Rooms {
	return &Rooms{saloons: make(map[string]*Room)}
}

// AddRoom creates a new room and add it to the set of rooms
// If the room already exists, do nothing
// Returns the roomID of the existing or new created room
func (rooms *Rooms) AddRoom(roomName string, owner *User, isPrivate bool) string {
	room := NewRoom(roomName, owner, isPrivate)

	_, exists := rooms.Room(room.ID)
	if exists == nil {
		return room.ID
	}

	rooms.mu.Lock()
	rooms.saloons[room.ID] = room
	rooms.mu.Unlock()

	return room.ID
}

// AddPeer add a peer to an existing room and sets roomID for the user
func (rooms *Rooms) AddPeer(roomID string, u *User) (*Room, error) {
	r, err := rooms.Room(roomID)
	if err != nil {
		return nil, err
	}

	err = r.AddPeer(u)
	if err != nil {
		return nil, err
	}

	u.RoomID = roomID

	return r, nil
}

// Peers return a user slice of connected peers in a room
func (rooms *Rooms) Peers(roomID string) (PeersType, error) {
	r, err := rooms.Room(roomID)
	if err != nil {
		return nil, err
	}

	return r.Peers, nil
}

// Rooms returns a room in a set of room given its ID
// Return an error if the room is not in set
func (rooms *Rooms) Room(roomID string) (*Room, error) {
	rooms.mu.RLock()
	r, ok := rooms.saloons[roomID]
	rooms.mu.RUnlock()

	if !ok {
		log.Println("access to unknown room, ID: " + roomID)
		return nil, errors.New(constants.ERR_ROOM_NIL)
	}
	return r, nil
}
