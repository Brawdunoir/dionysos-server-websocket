package requests

import (
	obj "github.com/Brawdunoir/dionysos-server/objects"
	responses "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

type CodeType string

// Requests codes
const (
	// Unregister a new user to the server and return the user's ID in the payload.
	DISCONNECTION = "DIS"
	// Register a new room to the server and return the room's ID in the payload.
	NEW_ROOM = "NRO"
	// Ask to join a room, return nothing. The answer is sent after owner decision.
	JOIN_ROOM = "JRO"
	// Follow a JOIN_ROOM request. Grant or deny user access to the room, return nothing.
	JOIN_ROOM_ANSWER = "JRA"
	// Forward the message to all peers within the room. The messages are not kept in the rooms.
	NEW_MESSAGE = "NME"
	// Change the username
	CHANGE_USERNAME = "CHU"
	// Change owner of room
	CHANGE_OWNER = "COW"
	// Quit a room
	QUIT_ROOM = "QRO"
	// Ask to kick a peer
	KICK_PEER = "KPE"
	// Ask to load a file in the room
	LOAD_FILE = "LFI"
	// Upload file chunk. It is only forwarded to peers and not keeped in the server
	UPLOAD_CHUNK = "UCH"
)

type IRequest interface {
	Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) responses.Response
	Code() CodeType
}
