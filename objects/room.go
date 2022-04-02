package objects

import (
	"errors"
	"sync"

	"github.com/Brawdunoir/dionysos-server/constants"
	"github.com/Brawdunoir/dionysos-server/utils"
	"go.uber.org/zap"
)

type PeersType []*User

// room represents data about a room for peers.
type Room struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	OwnerID   string     `json:"ownerid"`
	Peers     PeersType  `json:"peers"`
	IsPrivate bool       `json:"isPrivate"`
	mu        sync.Mutex `json:"-"`
}

func (r *Room) String() string {
	return r.Name + " (" + r.ID + ")"
}

// AddPeer safely adds a peer to a room and sets roomID for the user
func (r *Room) AddPeer(u *User, logger *zap.SugaredLogger) error {

	if ok := r.IsPeerPresent(u, logger); ok {
		logger.Errorw("add peer failed, user is already in the room", "user", u.ID, "username", u.Name, "room", r.ID, "roomname", r.Name)
		return errors.New(constants.ERR_PEER_ALREADY_IN_ROOM)
	}

	r.mu.Lock()
	r.Peers = append(r.Peers, u)
	r.mu.Unlock()

	u.RoomID = r.ID

	logger.Debugw("add peer", "user", u.ID, "username", u.Name, "room", r.ID, "roomname", r.Name)
	return nil
}

// SetOwnerID safely changes ownerID to room
func (r *Room) SetOwnerID(newOwnerID string) {
	r.mu.Lock()
	r.OwnerID = newOwnerID
	r.mu.Unlock()
}

// RemovePeer safely removes a peer from a room and sets roomID for the user
func (r *Room) RemovePeer(u *User, logger *zap.SugaredLogger) error {
	for i, p := range r.Peers {
		if p.ID == u.ID {
			r.mu.Lock()
			r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
			r.mu.Unlock()
			u.RoomID = ""
			logger.Debugw("remove peer from room", "user", u.ID, "username", u.Name, "room", r.ID, "roomname", r.Name)
			if u.ID == r.OwnerID && len(r.Peers) > 0 {
				newOwner := r.Peers[0]
				r.SetOwnerID(newOwner.ID)
				logger.Debugw("change owner of room", "oldOwner", u.ID, "oldOwnername", u.Name, "newOwner", newOwner.ID, "newOwnername", newOwner.Name, "room", r.ID, "roomname", r.Name)
			}
			return nil
		}
	}
	logger.Debugw("remove peer failed, the user cannot be found", "user", u.ID, "username", u.Name, "room", r.ID, "roomname", r.Name)
	return errors.New(constants.ERR_USER_NOT_IN_ROOM)
}

// IsPeerPresent evaluates if a certain user is in the room
func (r *Room) IsPeerPresent(u *User, logger *zap.SugaredLogger) bool {
	for _, p := range r.Peers {
		if p.ID == u.ID {
			return true
		}
	}
	return false
}

// SendJSONToPeers send the same json formatted message to all peers in the room
func (r *Room) SendJSONToPeers(json interface{}, logger *zap.SugaredLogger) {
	for _, peer := range r.Peers {
		peer.SendJSON(json, logger)
		logger.Debugw("peer has been notified of a message", "peer", peer.ID, "peername", peer.Name, "room", r.ID, "roomname", r.Name, "message", json)
	}
}

// Generate a room ID based on a roomname and an ownerPublicAddr
func generateRoomID(roomName, ownerPublicAddr string) string {
	return utils.GenerateStringHash(roomName + ownerPublicAddr)
}

// Creates a new room
func NewRoom(roomName string, owner *User, isPrivate bool) *Room {
	return &Room{ID: generateRoomID(roomName, owner.PublicIP), Name: roomName, OwnerID: owner.ID, Peers: PeersType{owner}, IsPrivate: isPrivate, mu: sync.Mutex{}}
}
