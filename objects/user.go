package objects

import (
	"sync"

	"github.com/gorilla/websocket"
)

// User defines a user
type user struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	RemoteAddr string          `json:"-"`
	Conn       *websocket.Conn `json:"-"`
	ConnMutex  sync.Mutex      `json:"-"`
}

func (u *user) String() string {
	return u.Name + " (" + u.RemoteAddr + ")"
}

// generateUserID generates an user ID based on a remote address and an username
func generateUserID(remoteAddr, username string) string {
	return generateStringHash(remoteAddr + username)
}

// newUser creates a new user
func newUser(username, remoteAddr string, conn *websocket.Conn) *user {
	return &user{ID: generateUserID(remoteAddr, username), Name: username, RemoteAddr: remoteAddr, ConnMutex: sync.Mutex{}, Conn: conn}
}
