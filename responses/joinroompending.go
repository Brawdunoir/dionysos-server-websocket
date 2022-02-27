package responses

import "encoding/json"

// JoinRoomPendingResponse is sent to the room's owner
// to ask permission for joining a room.
// RequesterID is set by the server and send to the room's owner.
// This way, the owner answer Yes or No to the request
// and join the RequesterID to his answer.
// See JoinRoomAnswerRequest
type JoinRoomPendingResponse struct {
	RoomID            string `json:"roomId"`
	RequesterUsername string `json:"requesterUsername"`
	RequesterID       string `json:"requesterId,omitempty"`
}

func (r JoinRoomPendingResponse) Code() CodeType {
	return JOIN_ROOM_PENDING
}

func (r JoinRoomPendingResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
