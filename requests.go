package main

const (
	NEWCONNECTION = "NCO"
	NEWROOM       = "NRO"
	JOINROOM      = "JRO"
)

type NewConnectionRequest struct {
	Username string `json:"username"`
}

type NewRoomRequest struct {
	RoomName  string `json:"roomname"`
	OwnerName string `json:"ownername"`
}
