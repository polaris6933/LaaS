package main

import (
	"LaaS/executor"
	"LaaS/life"
	"bufio"
	"crypto/sha256"
	"fmt"
	// "io/ioutil"
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

func (s *session) String() string {
	return fmt.Sprintf(s.getStringRepresentation())
}

func (s *session) authorize(user string) bool {
	return s.owner.name == user
}

func (s *Server) Register(username, password string) string {
	for _, user := range s.users {
		if user.name == username {
			return "user with the name " + username + " already exists"
		}
	}
	s.users = append(s.users, *newUser(username, password))
	return "registered user " + username
}

func (s *Server) Login(username, password string) string {
	for _, user := range s.users {
		if user.name == username {
			if user.password == sha256.Sum256([]byte(password)) {
				return "user " + username + " logged in"
			} else {
				return "invalid password for " + username
			}
		}
	}
	return "user " + username + " does not exists"
}

func (s *Server) Add(username, name string) string {
	for _, session := range s.sessions {
		if session.name == name {
			return "session with the name " + name + " already exists"
		}
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

func (s *Server) Remove(user, name string) string {
	// TODO: use sessionIndex()
	for idx, session := range s.sessions {
		if session.name == name {
			if !session.authorize(user) {
				return "user " + user + " not authorized"
			}
			last := len(s.sessions) - 1
			s.sessions[idx] = s.sessions[last]
			s.sessions = s.sessions[:last]
			return "removed session " + name
		}
	}
	return "no session with the name " + name + " found"
}

// TODO: authorization
func (s *Server) Start(user, session, config string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	s.sessions[index].currState = life.NewLife("predefined_configs/" + config)
	go func() {
		for {
			time.Sleep(time.Second)
			s.sessions[index].currState.NextGeneration()
		}
	}()

	return "successfully started session " + session
}

func (s *Server) Watch(user, session string) string {
	index := s.sessionIndex(session)
	if index == -1 {
		return "no session with the name " + session + " found"
	}
	return s.sessions[index].currState.Printable()
}

func pause(args []string, conection *net.Conn) {
	fmt.Println("pause()", args)
}

func resume(args []string, conection *net.Conn) {
	fmt.Println("resume()", args)
}

func export(args []string, conection *net.Conn) {
	fmt.Println("export()", args)
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
		fmt.Println("entered for")
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
		fmt.Println(connectionAddress, "-", response)
	}
	fmt.Println("this shit ends now")
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
