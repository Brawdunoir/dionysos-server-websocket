package main

type User struct {
	ID         string `json:"-"`
	Name       string `json:"name"`
	RemoteAddr string `json:"-"`
}

func (u User) String() string {
	return u.Name + " (" + u.RemoteAddr + ")"
}

// Generate an user ID based on a remote address and an username
func GenerateUserID(remoteAddr, username string) string {
	return generateStringHash(remoteAddr + username)
}

// Creates a new user
func NewUser(username, remoteAddr string) User {
	return User{ID: GenerateUserID(remoteAddr, username), Name: username, RemoteAddr: remoteAddr}
}
