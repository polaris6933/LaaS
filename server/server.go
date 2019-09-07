package main

import (
	"LaaS/executor"
	"LaaS/life"
	"bufio"
	"fmt"
	"net"
	"LaaS/server/session"
	"strings"
	"LaaS/server/user"
)

const connectionType = "tcp"
const port = ":8088"

type Server struct {
	sessions []*session.Session
	users    []user.User
}

func NewServer() *Server {
	s := new(Server)
	s.sessions = []*session.Session{}
	return s
}

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

func (s *Server) Register(username, password string) string {
	if s.userIndex(username) != -1 {
		return "user " + username + " already exists"
	}
	s.users = append(s.users, *user.NewUser(username, password))
	return "registered user " + username
}

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

func (s *Server) Add(username, name string) string {
	if s.sessionIndex(name) != -1 {
		return "session with the name " + name + " already exists"
	}
	owner := &s.users[s.userIndex(username)]
	s.sessions = append(s.sessions, session.NewSession(name, owner))
	return "successfully created session " + name
}

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

func (s *Server) Watch(_, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	return s.sessions[index].CurrState.Printable()
}

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
