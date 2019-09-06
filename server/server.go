package main

import (
	"LaaS/executor"
	"LaaS/life"
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"strings"
	"time"
)

const connectionType = "tcp"
const port = ":8088"
const timeFormat = "Mon Jan 2 2006 15:04"

type user struct {
	name     string
	password [sha256.Size]byte
}

func newUser(username, password string) *user {
	passwordHash := sha256.Sum256([]byte(password))
	return &user{name: username, password: passwordHash}
}

type session struct {
	owner     *user
	name      string
	created   time.Time
	currState *life.Life
	stopper   chan string
	isRunning bool
}

func NewSession(name string, owner *user) *session {
	s := new(session)
	s.name = name
	s.created = time.Now()
	s.owner = owner
	return s
}

type Server struct {
	sessions []*session
	users    []user
}

func NewServer() *Server {
	s := new(Server)
	s.sessions = []*session{}
	return s
}

func notAuthorized(name string) string {
	return "user " + name + " not authorized"
}

func noSession(name string) string {
	return "no session with the name " + name + " found"
}

func (s *Server) AssertExecutable() {}

func (s *session) getStringRepresentation() string {
	return "session " + s.name + ", created at " + s.created.Format(timeFormat)
}

func (s *Server) sessionIndex(name string) int {
	for idx, session := range s.sessions {
		if session.name == name {
			return idx
		}
	}
	return -1
}

func (s *session) run() {
	go func() {
		s.isRunning = true
		s.stopper = make(chan string)
		for {
			select {
			case <-s.stopper:
				s.isRunning = false
				close(s.stopper)
				return
			default:
				time.Sleep(time.Second)
				s.currState.NextGeneration()
			}
		}
	}()
}

func (s *session) stop() {
	s.stopper<- "stop"
}

func (s *session) String() string {
	return fmt.Sprintf(s.getStringRepresentation())
}

func (s *session) authorize(user string) bool {
	return s.owner.name == user
}

func (s *Server) userIndex(name string) int {
	for idx, user := range s.users {
		if user.name == name {
			return idx
		}
	}
	return -1
}

func (s *Server) Register(username, password string) string {
	if s.userIndex(username) != -1 {
		return "user with the name " + username + " already exists"
	}
	s.users = append(s.users, *newUser(username, password))
	return "registered user " + username
}

func (s *Server) Login(username, password string) string {
	index := s.userIndex(username)
	if index == -1 {
		return "user " + username + " does not exists"
	}
	user := s.users[index]
	if user.password == sha256.Sum256([]byte(password)) {
		return "user " + username + " logged in"
	} else {
		return "invalid password for " + username
	}
}

func (s *Server) Add(username, name string) string {
	if s.sessionIndex(name) != -1 {
		return "session with the name " + name + " already exists"
	}
	var owner *user
	for _, user := range s.users {
		if user.name == username {
			owner = &user
			break
		}
	}
	s.sessions = append(s.sessions, NewSession(name, owner))
	return "created session " + name
}

func (s *Server) Kill(user, name string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return noSession(name)
	}
	current := s.sessions[index]
	if !current.authorize(user) {
		return notAuthorized(user)
	}
	if current.isRunning {
		current.stop()
	}
	last := len(s.sessions) - 1
	s.sessions[index] = s.sessions[last]
	s.sessions = s.sessions[:last]
	return "session " + name + " successfuly killed"
}

func (s *Server) Start(user, name, config string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return noSession(name)
	}
	current := s.sessions[index]
	if !current.authorize(user) {
		return notAuthorized(user)
	}
	newLife, err := life.NewLife("predefined_configs/" + config)
	if err != nil {
		return err.Error()
	}

	current.currState = newLife
	current.run()

	return "successfully started session " + name
}

func (s *Server) Resume(user, name string) string {
	index := s.sessionIndex(name)
	if index == -1 {
		return noSession(name)
	}
	current := s.sessions[index]
	if !current.authorize(user) {
		return notAuthorized(user)
	}
	current.run()

	return "successfully resumed session " + name
}

func (s *Server) Stop(user, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	current := s.sessions[index]
	if !current.authorize(user) {
		return notAuthorized(user)
	}
	current.stop()
	return "session " + session + " successfully stopped"
}

func (s *Server) Watch(user, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	return s.sessions[index].currState.Printable()
}

func (s *Server) List() string {
	var listing strings.Builder
	for _, session := range s.sessions {
		listing.WriteString(session.getStringRepresentation())
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
