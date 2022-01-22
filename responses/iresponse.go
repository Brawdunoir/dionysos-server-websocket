package responses

// Responses status
const (
	// Basic success status
	SUCCESS = "SUC"
	// Basic error status
	ERROR   = "ERR"
	// Ask room's owner to take a decision on wathever to accept or deny room access to an user (requester)
	JOIN_ROOM_PENDING = "JRP"
	// Request had been accepted
	REQUEST_ACCEPTED = "REA"
	// Request had been refused
	REQUEST_DENIED = "RED"
	// Signal that a new peer joined the room
	NEW_PEER = "NEP"
)
