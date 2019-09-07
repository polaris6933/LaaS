// Package main provides the entry point of the server part of LaaS and
// contains most of the business logic of the application.
package main

import (
	"LaaS/executor"
	"LaaS/life"
	"LaaS/server/session"
	"LaaS/server/user"
	"bufio"
	"fmt"
	"net"
	"strings"
)

const connectionType = "tcp"
const port = ":8088"

// The Server is represented by a list of the users registered with the server
// and the sessions which the users have added.
// The Server methods return a human readable string which describes the result
// of the issued request - no matter if the operation has succeeded or failed.
type Server struct {
	sessions []*session.Session
	users    []user.User
}

// Constructs a Server.
func NewServer() *Server {
	s := new(Server)
	s.sessions = []*session.Session{}
	return s
}

// Implement the Executable interface for use with LaaS/executor.
func (s *Server) AssertExecutable() {}

func (s *Server) sessionIndex(name string) int {
	for idx, session := range s.sessions {
		if session.Name == name {
			return idx
		}
	}
	return -1
}

func (s *Server) userIndex(name string) int {
	for idx, user := range s.users {
		if user.Name == name {
			return idx
		}
	}
	return -1
}

// Registers a user with the server.
// Fails if the username is already taken.
func (s *Server) Register(username, password string) string {
	if s.userIndex(username) != -1 {
		return "user " + username + " already exists"
	}
	s.users = append(s.users, *user.NewUser(username, password))
	return "registered user " + username
}

// Logs the user in the server.
// Fails if
//   - a user with the name `username` does not exist
//   - the password `password` does not match the password of the user
func (s *Server) Login(username, password string) string {
	index := s.userIndex(username)
	if index == -1 {
		return "user " + username + " does not exist"
	}
	user := s.users[index]
	if user.Authorize(password) {
		return "user " + username + " logged in"
	} else {
		return "invalid password for " + username
	}
}

// Creates a new session whose owner is the user issuing the request.
// Fails if a sessions with the same name already exists.
func (s *Server) Add(username, name string) string {
	if s.sessionIndex(name) != -1 {
		return "session with the name " + name + " already exists"
	}
	owner := &s.users[s.userIndex(username)]
	s.sessions = append(s.sessions, session.NewSession(name, owner))
	return "successfully created session " + name
}

// Permanently removes a session from the server.
// Fails if
//   - no session with the name `name` exists
//   - the user issuing the request is not the owner of the session
func (s *Server) Kill(username, name string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return session.NoSession(name)
	}
	current := s.sessions[index]
	if !current.Authorize(username) {
		return user.NotAuthorized(username)
	}
	if current.IsRunning {
		current.Stop()
	}
	last := len(s.sessions) - 1
	s.sessions[index] = s.sessions[last]
	s.sessions = s.sessions[:last]
	return "session " + name + " successfully killed"
}

// Loads a game configuration in the session named `name` and starts the game.
// Fails if:
//   - a sessions with the name `name` does not exist
//   - the user issuing the request is not the owner of the session
//   - the session has already been started
//   - the config `config` does not exist
func (s *Server) Start(username, name, config string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return session.NoSession(name)
	}
	current := s.sessions[index]
	if !current.Authorize(username) {
		return user.NotAuthorized(username)
	}
	if current.IsRunning {
		return "session " + name + " is already running"
	}

	newLife, err := life.NewLife(config)
	if err != nil {
		return err.Error()
	}

	current.CurrState = newLife
	current.Run()

	return "successfully started session " + name
}

// Resumes a stopped session.
// Fails if:
//   - a sessions with the name `name` does not exist
//   - the user issuing the request is not the owner of the session
//   - the session is currently running
func (s *Server) Resume(username, name string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return session.NoSession(name)
	}
	current := s.sessions[index]
	if !current.Authorize(username) {
		return user.NotAuthorized(username)
	}
	if current.IsRunning {
		return "session " + name + " is already running"
	}
	current.Run()

	return "successfully resumed session " + name
}

// Temporarily stops a running session.
// Fails if:
//   - a sessions with the name `name` does not exist
//   - the user issuing the request is not the owner of the session
//   - the session is not currently running
func (s *Server) Stop(username, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	current := s.sessions[index]
	if !current.Authorize(username) {
		return user.NotAuthorized(username)
	}
	if !current.IsRunning {
		return "session " + session + " is already stopped"
	}
	current.Stop()
	return "session " + session + " successfully stopped"
}

// Returns the current state of the running game associated with the session
// named `session`. Any user can watch any session.
// Fails if:
//   - a sessions with the name `name` does not exist
//   - the session had not been started
func (s *Server) Watch(_, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	if s.sessions[index].CurrState == nil {
		return "the session " + session + " has not been started"
	}
	return s.sessions[index].CurrState.Printable()
}

// Returns information for all the sessions on the server.
func (s *Server) List() string {
	var listing strings.Builder
	for _, session := range s.sessions {
		listing.WriteString(session.GetStringRepresentation())
		listing.WriteString("\n")
	}
	listing.WriteString(fmt.Sprintf("\n%d total", len(s.sessions)))
	return listing.String()
}

func (s *Server) handleRequest(connection *net.Conn) {
	connectionAddress := (*connection).RemoteAddr().String()
	fmt.Println("serving", connectionAddress)
	var response string
	for {
		received, err := bufio.NewReader(*connection).ReadString('\000')
		if err != nil {
			fmt.Println(err)
			return
		}
		request := strings.TrimRight(string(received), "\000")
		execResult, err := executor.Execute(s, request)
		if err != nil {
			fmt.Println(err)
			response = "internal server error"
		} else {
			response = execResult[0].Interface().(string)
		}
		(*connection).Write([]byte(response + "\000"))
		if !strings.HasPrefix(request, "watch") {
			fmt.Println(connectionAddress, "-", response)
		}
	}
}

func main() {
	l, err := net.Listen(connectionType, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	s := NewServer()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()
		go s.handleRequest(&c)
	}
}
