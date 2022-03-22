package objects

import (
	"sync"

	"github.com/Brawdunoir/dionysos-server/utils"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// User defines a user.
// Salt is used to generate the ID, it cannot change during a session
// and it is a secret shared between the client and the server.
type User struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	RoomID    string          `json:"-"`
	PublicIP  string          `json:"-"`
	Conn      *websocket.Conn `json:"-"`
	ConnMutex sync.Mutex      `json:"-"`
}

func (u *User) String() string {
	return u.Name + " (" + u.PublicIP + ")"
}

// SendJSON send a json formatted message to a user, respecting concurrency
func (u *User) SendJSON(json interface{}, logger *zap.SugaredLogger) {
	u.ConnMutex.Lock()
	err := u.Conn.WriteJSON(json)
	u.ConnMutex.Unlock()
	if err != nil {
		logger.Errorw("send json failed", "user", u.ID, "username", u.Name)
	}
}

// generateUserID generates an user ID based on a public address and a salt send by the client
func generateUserID(publicAddr, salt string) string {
	return utils.GenerateStringHash(publicAddr + salt)
}

// newUser creates a new user
func NewUser(username, publicAddr, salt string, conn *websocket.Conn) *User {
	return &User{ID: generateUserID(publicAddr, salt), RoomID: "", Name: username, PublicIP: publicAddr, ConnMutex: sync.Mutex{}, Conn: conn}
}
