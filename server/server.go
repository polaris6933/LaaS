package main

import (
	"LaaS/executor"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func add(args []string, conection *net.Conn) {
	fmt.Println("add()", args)
}

func remove(args []string, conection *net.Conn) {
	fmt.Println("remove()", args)
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

func list(args []string, conection *net.Conn) {
	fmt.Println("list()", args)
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
		"add":    {add, 2},
		"remove": {remove, 1},
		"pause":  {pause, 1},
		"resume": {resume, 1},
		"export": {export, 2},
		"list":   {list, 0},
	}

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
