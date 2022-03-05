package objects

import (
	"sync"

	"github.com/gorilla/websocket"
)

// User defines a user.
// Salt is used to generate the ID, it cannot change during a session
// and it is a secret shared between the client and the server.
type User struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Salt       string          `json:"-"`
	RemoteAddr string          `json:"-"`
	Conn       *websocket.Conn `json:"-"`
	ConnMutex  sync.Mutex      `json:"-"`
}

func (u *User) String() string {
	return u.Name + " (" + u.RemoteAddr + ")"
}

// generateUserID generates an user ID based on a remote address and a salt send by the client
func generateUserID(remoteAddr, salt string) string {
	return generateStringHash(remoteAddr + salt)
}

// newUser creates a new user
func NewUser(username, remoteAddr, salt string, conn *websocket.Conn) *User {
	return &User{ID: generateUserID(remoteAddr, salt), Salt: salt, Name: username, RemoteAddr: remoteAddr, ConnMutex: sync.Mutex{}, Conn: conn}
}
