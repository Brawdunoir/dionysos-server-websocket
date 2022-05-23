package responses

type CodeType string

// Response possible codes
const (
	// Basic error response
	ERROR = "ERR"
	// New connection response
	CONNECTION_SUCCESS = "COS"
	// New room response
	ROOM_CREATION_SUCCESSS = "RCS"
	// Send basic room information while we ping the owner for confirmation or we add the peer for public room
	JOIN_ROOM = "JRO"
	// Ask room's owner to take a decision on wathever to accept or deny room access to an user (requester)
	JOIN_ROOM_PENDING = "JRP"
	// Request has been denied
	DENIED = "DEN"
	// Signal that a new peer joined the room
	NEW_PEER = "NEP"
	// Signal that a new message has been sent with the content
	NEW_MESSAGE = "NME"
	// Signal that the username has been changed internally
	CHANGE_USERNAME = "CHU"
	// Signal that a new file has been loaded (and that fragments will be sent)
	LOAD_FILE = "LFI"
	// Upload file chunk. It is only forwarded to peers and not keeped in the server
	UPLOAD_CHUNK = "UCH"
)

type IResponse interface {
	Code() CodeType
	Marshal() ([]byte, error)
}
