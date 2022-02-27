package responses

type CodeType string

// Response possible codes
const (
	// Basic success response
	SUCCESS = "SUC"
	// Basic error response
	ERROR = "ERR"
	// New connection response
	CONNECTION_SUCCESS = "COS"
	// New room response
	ROOM_CREATION_SUCCESSS = "RCS"
	// Ask room's owner to take a decision on wathever to accept or deny room access to an user (requester)
	JOIN_ROOM_PENDING = "JRP"
	// Request has been denied
	DENIED = "DEN"
	// Signal that a new peer joined the room
	NEW_PEER = "NEP"
)

type IResponse interface {
	Code() CodeType
	Marshal() ([]byte, error)
}
