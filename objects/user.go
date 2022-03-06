package objects

import (
	"sync"

	"github.com/Brawdunoir/dionysos-server/utils"
	"github.com/gorilla/websocket"
)

// User defines a user.
// Salt is used to generate the ID, it cannot change during a session
// and it is a secret shared between the client and the server.
type User struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	RoomID    string          `json:"-"`
	Salt      string          `json:"-"`
	PublicIP  string          `json:"-"`
	Conn      *websocket.Conn `json:"-"`
	ConnMutex sync.Mutex      `json:"-"`
}

func (u *User) String() string {
	return u.Name + " (" + u.PublicIP + ")"
}

// generateUserID generates an user ID based on a public address and a salt send by the client
func generateUserID(publicAddr, salt string) string {
	return utils.GenerateStringHash(publicAddr + salt)
}

// newUser creates a new user
func NewUser(username, publicAddr, salt string, conn *websocket.Conn) *User {
	return &User{ID: generateUserID(publicAddr, salt), RoomID: "", Salt: salt, Name: username, PublicIP: publicAddr, ConnMutex: sync.Mutex{}, Conn: conn}
}
