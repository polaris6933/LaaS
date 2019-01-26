package main

import (
	"LaaS/executor"
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	// "strconv"
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
	sessions []*session
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

func (s *session) String() string {
	return fmt.Sprintf(s.getStringRepresentation())
}

func (s *Server) Add(name, pass string) string {
	for _, session := range s.sessions {
		if session.name == name {
			return "session with the name " + name + " already exists"
		}
	}
	s.sessions = append(s.sessions, NewSession(name, pass))
	return "created session " + name
}

func (s *Server) Remove(name string) string {
	for idx, session := range s.sessions {
		if session.name == name {
			last := len(s.sessions) - 1
			s.sessions[idx] = s.sessions[last]
			s.sessions = s.sessions[:last]
			return "removed session " + name
		}
	}
	return "no session with the name" + name + "found"
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
		received, err := bufio.NewReader(*connection).ReadString('\000')
		if err != nil {
			fmt.Println(err)
			return
		}
		request := strings.TrimRight(string(received), "\000")
		fmt.Println(len(request))
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
		go s.handleRequest(&c)
	}
}
