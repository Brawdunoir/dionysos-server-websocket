package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	RemoteAddr    string          `json:"-"`
	Accepted      chan bool       `json:"-"`
	AcceptedMutex sync.Mutex      `json:"-"`
	Conn          *websocket.Conn `json:"-"`
	ConnMutex     sync.Mutex      `json:"-"`
}

func (u *User) String() string {
	return u.Name + " (" + u.RemoteAddr + ")"
}

// Generate an user ID based on a remote address and an username
func GenerateUserID(remoteAddr, username string) string {
	return generateStringHash(remoteAddr + username)
}

// Creates a new user
func NewUser(username, remoteAddr string, conn *websocket.Conn) *User {
	return &User{ID: GenerateUserID(remoteAddr, username), Name: username, RemoteAddr: remoteAddr, ConnMutex: sync.Mutex{}, Conn: conn, Accepted: make(chan bool)}
}
