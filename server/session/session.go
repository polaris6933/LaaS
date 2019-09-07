// Package session contains the `Session` type which represents a game session
// in the server.
package session

import (
	"fmt"
	"time"
	"LaaS/server/user"
	"LaaS/life"
)

const timeFormat = "Mon Jan 2 2006 15:04"

// A session is represented by its name, owner, creation time and the current
// state of the game. It also contains a channel used to signal the game to
// stop and a flag indicating if the game is currently running or not.
type Session struct {
	owner     *user.User
	Name      string
	created   time.Time
	CurrState *life.Life
	stopper   chan struct{}
	IsRunning bool
}

// Constructs a new session.
func NewSession(name string, owner *user.User) *Session {
	s := new(Session)
	s.Name = name
	s.created = time.Now()
	s.owner = owner
	return s
}

// Returns a string describing the session (human readable).
func (s *Session) GetStringRepresentation() string {
	return "session " + s.Name + ", created at " + s.created.Format(timeFormat)
}

// Begins iteration the generations of the game.
// DO NOT call Run() on sessions whose state has not been initialized.
func (s *Session) Run() {
	go func() {
		s.IsRunning = true
		s.stopper = make(chan struct{})
		for {
			select {
			case <-s.stopper:
				s.IsRunning = false
				close(s.stopper)
				return
			default:
				time.Sleep(time.Second)
				s.CurrState.NextGeneration()
			}
		}
	}()
}

// Signals the session to stop executing the game.
func (s *Session) Stop() {
	s.stopper<- struct{}{}
}

// Implement the Stringer interface.
func (s *Session) String() string {
	return fmt.Sprintf(s.GetStringRepresentation())
}

// Checks whether the given `user` is the owner the session or not.
func (s *Session) Authorize(user string) bool {
	return s.owner.Name == user
}

// Returns a message (string) stating that the session named `name` is
// non-existent.
func NoSession(name string) string {
	return "no session with the name " + name + " found"
}
