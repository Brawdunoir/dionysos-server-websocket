package responses

import "encoding/json"

// JoinRoomResponse correspond to a successful JoinRoomRequest.
// It sends basic room information such as roomname, roomid, and
// isPrivate.
type JoinRoomResponse struct {
	RoomName  string `json:"roomName"`
	RoomID    string `json:"roomId"`
	IsPrivate bool   `json:"isPrivate"`
}

func (r JoinRoomResponse) Code() CodeType {
	return JOIN_ROOM
}

func (r JoinRoomResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
