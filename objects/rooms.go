package objects

import (
	"errors"
	"sync"

	"github.com/Brawdunoir/dionysos-server/constants"
	"go.uber.org/zap"
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
func (rooms *Rooms) AddRoom(roomName string, owner *User, isPrivate bool, logger *zap.SugaredLogger) string {
	room := NewRoom(roomName, owner, isPrivate)

	if _, exists := rooms.secureRoom(room.ID); exists {
		logger.Debugw("add room, room already exists", "room", room.ID, "roomname", roomName, "owner", owner.ID, "ownername", owner.Name)
		return room.ID
	}

	rooms.mu.Lock()
	rooms.saloons[room.ID] = room
	rooms.mu.Unlock()

	logger.Debugw("add room", "room", room.ID, "roomname", roomName, "owner", owner.ID, "ownername", owner.Name)
	return room.ID
}

// AddPeer add a peer to an existing room and sets roomID for the user
func (rooms *Rooms) AddPeer(roomID string, u *User, logger *zap.SugaredLogger) (*Room, error) {
	r, err := rooms.Room(roomID, logger)
	if err != nil {
		return nil, err
	}

	err = r.AddPeer(u, logger)
	if err != nil {
		return nil, err
	}

	u.RoomID = roomID

	return r, nil
}

// Peers return a user slice of connected peers in a room
func (rooms *Rooms) Peers(roomID string, logger *zap.SugaredLogger) (PeersType, error) {
	r, err := rooms.Room(roomID, logger)
	if err != nil {
		return nil, err
	}

	return r.Peers, nil
}

// Rooms returns a room in a set of room given its ID
// Return an error if the room is not in set
func (rooms *Rooms) Room(roomID string, logger *zap.SugaredLogger) (*Room, error) {
	r, ok := rooms.secureRoom(roomID)

	if !ok {
		logger.Debugw("room does not exist", "room", roomID)
		return nil, errors.New(constants.ERR_ROOM_NIL)
	}
	return r, nil
}

func (rooms *Rooms) secureRoom(roomID string) (r *Room, ok bool) {
	rooms.mu.RLock()
	r, ok = rooms.saloons[roomID]
	rooms.mu.RUnlock()
	return
}
