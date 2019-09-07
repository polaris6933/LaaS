package user

import "crypto/sha256"

type User struct {
	Name     string
	password [sha256.Size]byte
}

func NewUser(username, password string) *User {
	passwordHash := sha256.Sum256([]byte(password))
	return &User{Name: username, password: passwordHash}
}

func (u *User) Authorize(password string) bool {
	return u.password == sha256.Sum256([]byte(password))
}

func NotAuthorized(name string) string {
	return "user " + name + " not authorized"
}
