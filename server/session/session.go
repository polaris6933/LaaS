package session

import (
	"fmt"
	"time"
	"LaaS/server/user"
	"LaaS/life"
)

const timeFormat = "Mon Jan 2 2006 15:04"

type Session struct {
	owner     *user.User
	Name      string
	created   time.Time
	CurrState *life.Life
	// TODO: use struct{}
	stopper   chan string
	IsRunning bool
}

func NewSession(name string, owner *user.User) *Session {
	s := new(Session)
	s.Name = name
	s.created = time.Now()
	s.owner = owner
	return s
}

func (s *Session) GetStringRepresentation() string {
	return "session " + s.Name + ", created at " + s.created.Format(timeFormat)
}

func (s *Session) Run() {
	go func() {
		s.IsRunning = true
		s.stopper = make(chan string)
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

func (s *Session) Stop() {
	s.stopper<- "stop"
}

func (s *Session) String() string {
	return fmt.Sprintf(s.GetStringRepresentation())
}

func (s *Session) Authorize(user string) bool {
	return s.owner.Name == user
}

func NoSession(name string) string {
	return "no session with the name " + name + " found"
}
