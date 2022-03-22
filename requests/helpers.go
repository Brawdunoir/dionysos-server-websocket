package requests

import (
	obj "github.com/Brawdunoir/dionysos-server/objects"
	"go.uber.org/zap"
)

func getUserAndRoom(userSalt, publicAddr, roomID string, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (user *obj.User, room *obj.Room, err error) {
	user, err = users.User(userSalt, publicAddr, logger)
	if err != nil {
		return
	}

	room, err = rooms.Room(roomID, logger)
	if err != nil {
		return
	}

	return
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

func getUserAndRoomAndRoomOwner(userSalt, publicAddr, roomID string, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (user, owner *obj.User, room *obj.Room, err error) {
	user, room, err = getUserAndRoom(userSalt, publicAddr, roomID, users, rooms, logger)
	if err != nil {
		return
	}

	owner, err = users.UserByID(room.OwnerID, logger)
	if err != nil {
		return
	}

	return
}