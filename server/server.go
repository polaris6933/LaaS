package main

import (
	"LaaS/executor"
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const timeFormat string = "Mon Jan 2 2006 15:04"

type session struct {
	name         string
	created      time.Time
	passwordHash [sha256.Size]byte
}

func NewSession(name string, password string) *session {
	s := new(session)
	s.name = name
	s.created = time.Now()
	s.passwordHash = sha256.Sum256([]byte(password))
	return s
}

type Server struct {
	sessions   []*session
	connection *net.Conn
}

func NewServer() *Server {
	s := new(Server)
	s.sessions = []*session{}
	s.connection = nil
	return s
}

func (s *Server) AssertExecutable() {}

func (s *session) getStringRepresentation() string {
	return "session " + s.name + ", created at " + s.created.Format(timeFormat)
}

func (s *session) String() string {
	return fmt.Sprintf(s.getStringRepresentation())
}

func (s *Server) Add(name, pass string) {
	for _, session := range s.sessions {
		if session.name == name {
			(*s.connection).Write(
				[]byte("session with the name " + name + "already exists"))
			return
		}
	}
	s.sessions = append(s.sessions, NewSession(name, pass))
	fmt.Println((*s.connection).RemoteAddr().String(), "created session", name)
}

func (s *Server) Remove(name string) {
	for idx, session := range s.sessions {
		if session.name == name {
			last := len(s.sessions) - 1
			s.sessions[idx] = s.sessions[last]
			s.sessions = s.sessions[:last]
			return
		}
	}
	fmt.Println("no session with the name", name, "found")
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

func (s *Server) list() {
	for _, session := range s.sessions {
		fmt.Println(session.getStringRepresentation())
	}
	fmt.Println(len(s.sessions), "total\n")
}

func (s *Server) handleRequest() {
	c := *s.connection
	fmt.Println("serving", c.RemoteAddr().String())
	for {
		received, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		request := strings.TrimSpace(string(received))
		executor.Execute(s, request)
		s.list()
	}
}

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("argument error")
		return
	}

	connectionType := args[1]
	port := ":" + args[2]
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
		s.connection = &c
		go s.handleRequest()
	}
}
