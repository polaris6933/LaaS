// Package user contains the `User` type which represents a registered user
// in the server.
package user

import "crypto/sha256"

// A user is represented by their name and password.
type User struct {
	Name     string
	password [sha256.Size]byte
}

// Constructs a new user.
func NewUser(username, password string) *User {
	passwordHash := sha256.Sum256([]byte(password))
	return &User{Name: username, password: passwordHash}
}

// Checks wether the given `password` hash matched that of the user.
func (u *User) Authorize(password string) bool {
	return u.password == sha256.Sum256([]byte(password))
}

// Returns a message (string) stating that the user is not authorized for some
// action.
func NotAuthorized(name string) string {
	return "user " + name + " not authorized"
}
