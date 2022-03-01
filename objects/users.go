package objects

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Users struct {
	members map[string]*User
	mu      sync.RWMutex
}

// Creates a new set of Users
func NewUsers() *Users {
	return &Users{members: make(map[string]*User)}
}

// AddUser creates a new user and add it to the set of users
// If the user already exists, do nothing
// Returns user ID
func (users *Users) AddUser(username, remoteAddr string, conn *websocket.Conn) string {
	if u, exists := users.User(username, remoteAddr); exists == nil {
		return u.ID
	}

	user := NewUser(username, remoteAddr, conn)

	users.mu.Lock()
	users.members[user.ID] = user
	users.mu.Unlock()

	return user.ID
}

// UserByID returns a user in a set of user given its ID
// Return an error if the user is not in set
func (users *Users) UserByID(userID string) (*User, error) {
	users.mu.RLock()
	u, exists := users.members[userID]
	users.mu.RUnlock()

	if !exists {
		log.Println("access to unknown user, ID: " + userID)
		return nil, errors.New(USER_NIL)
	} else {
		return u, nil
	}
}

// User returns a user in a set of user given its username and remote address
// Return an error if the user is not in set
func (users *Users) User(username, remoteAddr string) (*User, error) {
	userID := generateUserID(remoteAddr, username)

	return users.UserByID(userID)
}
