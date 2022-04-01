package requests

import (
	"errors"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// Handle a disconnection from a client.
func DisconnectPeer(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	err := removeUserFromRoom(client, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}
	users.RemoveUser(client.ID, logger)

	logger.Infow("disconnection", "user", client.ID, "username", client.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: DISCONNECTION}, logger)
	return
}

// removeUserFromRoom remove the user from the room and notify the remaining peers of the changement.
func removeUserFromRoom(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) error {
	if client.RoomID == "" {
		logger.Debugw("not part of any room", "user", client.ID, "username", client.Name)
	}
	room, err := rooms.RemovePeer(client.RoomID, client, logger)
	if err != nil {
		return err
	}

	notifyPeers(rooms, room, logger)

	return nil
}

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
// Take care of empty or nil room
func notifyPeers(rooms *obj.Rooms, room *obj.Room, logger *zap.SugaredLogger) error {
	if room == nil || len(room.Peers) == 0 {
		logger.Debug("room is empty or does not exists")
		return errors.New(constants.ERR_ROOM_NIL)
	}
	peers, err := rooms.Peers(room.ID, logger)
	if err != nil {
		return err
	}

	// Send updated peers list to all peers
	mes := res.NewResponse(res.NewPeersResponse{Peers: peers, OwnerID: room.OwnerID}, logger)
	room.SendJSONToPeers(mes, logger)

	logger.Infow("notify peers", "room", room.ID, "roomname", room.Name)

	return nil
}
