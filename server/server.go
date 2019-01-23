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

var sessions []*session

const timeFormat string = "Mon Jan 2 2006 15:04"

type session struct {
	name                  string
	created               time.Time
	passwordHash          [sha256.Size]byte
}

func NewSession(name string, password string) *session {
	s := new(session)
	s.name = name
	s.created = time.Now()
	s.passwordHash = sha256.Sum256([]byte(password))
	return s
}

func (s *session) getStringRepresentation() string {
	return "session " + s.name + ", created at " + s.created.Format(timeFormat)
}

func (s *session) String() string {
	return fmt.Sprintf(s.getStringRepresentation())
}

func add(args []string, connection *net.Conn) {
	name := args[0]
	pass := args[1]
	for _, session := range sessions {
		if session.name == name {
			(*connection).Write(
				[]byte("session with the name " + name + "already exists"))
			return
		}
	}
	sessions = append(sessions, NewSession(name, pass))
	fmt.Println((*connection).RemoteAddr().String(), "created session", name)
}

func remove(args []string, conection *net.Conn) {
	name := args[0]
	for idx, session := range sessions {
		if session.name == name {
			last := len(sessions) - 1
			sessions[idx] = sessions[last]
			sessions = sessions[:last]
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

func list(args []string, connection *net.Conn) {
	for _, session := range sessions {
		(*connection).Write([]byte(session.getStringRepresentation()))
	}
	fmt.Println(len(sessions), "total\n")
}

func handleRequest(e *executor.Executor) {
	c := *e.Connection
	fmt.Println("serving", c.RemoteAddr().String())
	for {
		received, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		request := strings.TrimSpace(string(received))
		e.Execute(request)
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

	var actions = map[string]executor.Action{
		"add":          {add, 2},
		"remove":       {remove, 1},
		"pause":        {pause, 1},
		"resume":       {resume, 1},
		"export":       {export, 2},
		"list":         {list, 0},
		"authenticate": {authenticate, 1},
	}

	sessions = []*session{}
	executor := executor.NewExecutor(actions)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()
		executor.Connection = &c
		go handleRequest(executor)
	}
}
