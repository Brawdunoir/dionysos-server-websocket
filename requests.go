package main

import (
	"errors"
	"fmt"
)

const (
	NEWCONNECTION = "NCO"
	NEWROOM       = "NRO"
	JOINROOM      = "JRO"
)

type NewConnectionRequest struct {
	Username string `json:"username"`
}

func checkNewConnectionRequest(r NewConnectionRequest) error {
	if r.Username == "" {
		return errors.New("username is empty")
	}

	return nil
}

type NewRoomRequest struct {
	RoomName  string `json:"roomname"`
	OwnerName string `json:"ownername"`
}

func checkNewRoomRequest(r NewRoomRequest) error {
	var err error

	if r.RoomName == "" {
		err = errors.New("roomname is empty")
	}
	if r.OwnerName == "" {
		err = fmt.Errorf("%w; ownername is empty", err)
	}

	return err
}
