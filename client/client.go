package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type clientAction struct {
	call   func([]string, *net.Conn)
	argCnt int
}

type executor struct {
	actions    map[string]clientAction
	connection net.Conn
}

func NewExecutor(actions map[string]clientAction) *executor {
	var e executor
	e.actions = actions
	e.connection = nil
	return &e
}

func (e *executor) closer() {
	disconnect([]string{}, &e.connection)
}

func (e *executor) execute(command string) {
	commandSplit := strings.Split(command, " ")
	commandName := commandSplit[0]
	commandArgs := commandSplit[1:]

	action, ok := e.actions[commandName]
	if !ok {
		fmt.Println(commandName, "is not a valid operation")
		return
	}
	if len(commandArgs) != action.argCnt {
		fmt.Println("wrong number of arguments for", commandName)
		return
	}
	action.call(commandArgs, &e.connection)
}

func connect(args []string, connection *net.Conn) {
	fmt.Println("connect()", args)
}

func disconnect(args []string, connection *net.Conn) {
	fmt.Println("disconnect()", args)
}

func start(args []string, connection *net.Conn) {
	fmt.Println("start()", args)
}

func kill(args []string, connection *net.Conn) {
	fmt.Println("kill()", args)
}

func unlock(args []string, connection *net.Conn) {
	fmt.Println("unlock()", args)
}

func save(args []string, connection *net.Conn) {
	fmt.Println("save()", args)
}

func pause(args []string, connection *net.Conn) {
	fmt.Println("pause()", args)
}

func resume(args []string, connection *net.Conn) {
	fmt.Println("resume()", args)
}

func exit(args []string, connection *net.Conn) {
	disconnect([]string{}, connection)
	os.Exit(0)
}

func main() {
	var actions = map[string]clientAction{
		"connect":    {connect, 2},
		"disconnect": {disconnect, 0},
		"start":      {start, 2},
		"kill":       {kill, 1},
		"unlock":     {unlock, 1},
		"save":       {save, 2},
		"pause":      {pause, 1},
		"resume":     {resume, 1},
		"exit":       {exit, 0},
	}
	executor := NewExecutor(actions)
	defer executor.closer()
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		fmt.Println(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		executor.execute(strings.TrimSuffix(input, "\n"))
	}
}
