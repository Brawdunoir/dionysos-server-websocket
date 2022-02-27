package responses

import "encoding/json"

type NewRoomResponse struct {
	RoomID string `json:"roomId"`
}

func (r NewRoomResponse) Code() CodeType {
	return ROOM_CREATION_SUCCESSS
}

func (r NewRoomResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
