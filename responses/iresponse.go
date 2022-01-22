package responses

// Responses status
const (
	SUCCESS = "SUC"
	ERROR   = "ERR"
	// Ask room's owner to take a decision on wathever to accept or deny room access to an user (requester)
	JOIN_ROOM_PENDING = "JRP"
	// Request had been accepted
	REQUEST_ACCEPTED = "REA"
	// Request had been refused
	REQUEST_DENIED = "RED"
)
