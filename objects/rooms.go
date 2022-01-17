package objects

import (
	"errors"
	"sync"
)

type Rooms struct {
	saloons map[string]*room
	mu      sync.RWMutex
}

// Creates a new set of Rooms
func NewRooms() *Rooms {
	return &Rooms{saloons: make(map[string]*room)}
}

// AddRoom creates a new room and add it to the set of rooms
// If the room already exists, do nothing
// Returns the roomID of the existing or new created room
func (rooms *Rooms) AddRoom(roomName string, owner *user) string {
	room := newRoom(roomName, owner)

	_, exists := rooms.Room(room.ID)
	if exists == nil {
		return room.ID
	}

	rooms.mu.Lock()
	rooms.saloons[room.ID] = room
	rooms.mu.Unlock()

	return room.ID
}

// Rooms returns a room in a set of room given its ID
// Return an error if the room is not in set
func (rooms *Rooms) Room(roomID string) (*room, error) {
	rooms.mu.RLock()
	r, ok := rooms.saloons[roomID]
	rooms.mu.RUnlock()

	if !ok {
		return nil, errors.New("room does not exist, ID: " + roomID)
	}
	return r, nil
}
