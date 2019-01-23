package main

import (
	"LaaS/executor"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func connect(args []string, connection *net.Conn) {
	connectionType := args[0]
	connectTo := args[1]
	disconnect([]string{}, connection)
	var err error
	*connection, err = net.Dial(connectionType, connectTo)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func disconnect(args []string, connection *net.Conn) {
	if *connection == nil {
		return
	}
	(*connection).Close()
	*connection = nil
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

func send(args []string, connection *net.Conn) {
	fmt.Fprintf(*connection, strings.Join(args, " ")+"\n")
}

func exit(args []string, connection *net.Conn) {
	disconnect([]string{}, connection)
	os.Exit(0)
}

func main() {
	var actions = map[string]executor.Action{
		"connect":    {connect, 2},
		"disconnect": {disconnect, 0},
		"start":      {start, 2},
		"kill":       {kill, 1},
		"unlock":     {unlock, 1},
		"save":       {save, 2},
		"pause":      {pause, 1},
		"resume":     {resume, 1},
		"exit":       {exit, 0},
		"send":       {send, 1},
	}
	executor := executor.NewExecutor(actions)
	defer disconnect([]string{}, executor.Connection)
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		executor.Execute(strings.TrimSuffix(input, "\n"))
	}
}
