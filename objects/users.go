package objects

import (
	"errors"
	"sync"

	"github.com/Brawdunoir/dionysos-server/constants"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
func (users *Users) AddUser(username, publicAddr, salt string, conn *websocket.Conn, logger *zap.SugaredLogger) *User {
	user := NewUser(username, publicAddr, salt, conn)

	if _, exists := users.secureUserByID(user.ID); exists {
		logger.Debugw("add user, user already exists", "user", user.ID, "username", user.Name)
		return user
	}

	users.mu.Lock()
	users.members[user.ID] = user
	users.mu.Unlock()

	logger.Debugw("add user", "user", user.ID, "username", user.Name)

	return user
}

func (users *Users) ChangeUsername(userID, newUsername string, logger *zap.SugaredLogger) error {
	user, err := users.UserByID(userID, logger)
	if err != nil {
		logger.Debugw("change username failed", "user", user.ID, "username", user.Name)
		return err
	}

	users.mu.Lock()
	oldName := user.Name
	user.Name = newUsername
	users.mu.Unlock()

	logger.Debugw("change username", "user", user.ID, "new username", user.Name, "old username", oldName)

	return nil
}

// UserByID returns a user in a set of user given its ID
// Return an error if the user is not in set
func (users *Users) UserByID(userID string, logger *zap.SugaredLogger) (*User, error) {
	u, exists := users.secureUserByID(userID)

	if !exists {
		logger.Debugw("user does not exist", "user", userID)
		return nil, errors.New(constants.ERR_USER_NIL)
	} else {
		return u, nil
	}
}

// User returns a user in a set of user given its salt and public address
// Return an error if the user is not in set
func (users *Users) User(salt, publicAddr string, logger *zap.SugaredLogger) (*User, error) {
	userID := generateUserID(publicAddr, salt)
	return users.UserByID(userID, logger)
}

func (users *Users) secureUserByID(userID string) (u *User, ok bool) {
	users.mu.RLock()
	u, ok = users.members[userID]
	users.mu.RUnlock()
	return
}
