package requests

import (
	"errors"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

func getUserByIdAndRoom(userID, roomID string, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (user *obj.User, room *obj.Room, err error) {
	user, err = users.UserByID(userID, logger)
	if err != nil {
		return
	}

	room, err = rooms.Room(roomID, logger)
	if err != nil {
		return
	}

	return
}

func getRoomAndRoomOwner(roomID string, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (room *obj.Room, owner *obj.User, err error) {

	room, err = rooms.Room(roomID, logger)
	if err != nil {
		return
	}

	owner, err = users.UserByID(room.OwnerID, logger)
	if err != nil {
		return
	}

	return
}

func addPeerAndNotify(requester *obj.User, rooms *obj.Rooms, room *obj.Room, logger *zap.SugaredLogger) error {
	// Add the newcoming to the list of the peer before notifying
	_, err := rooms.AddPeer(room.ID, requester, logger)
	if err != nil {
		return err
	}

	err = notifyPeers(rooms, room, logger)
	if err != nil {
		return err
	}

	return nil
}

// Notify peers that the room peer list has changed
func notifyPeers(rooms *obj.Rooms, room *obj.Room, logger *zap.SugaredLogger) error {
	peers, err := rooms.Peers(room.ID, logger)
	if err != nil {
		return errors.New("error when retrieving peers in room")
	}

	// Send updated peers list to all peers
	mes := res.NewResponse(res.NewPeersResponse{Peers: peers, OwnerID: room.OwnerID}, logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("notify peers", "room", room.ID, "roomname", room.Name)

	return nil
}
