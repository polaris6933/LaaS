package main

import (
	"LaaS/executor"
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
	"strings"
	"syscall"
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
	if *connection == nil {
		fmt.Println("not connected to server atm")
		return
	}

	name := args[0]
	password := passwordConfirmation()
	fmt.Fprintf(*connection, "add "+name+" "+password+"\n")
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

func passwordConfirmation() string {
	var firstAttempt, secondAttempt string
	for {
		firstAttempt = readPassword("input password: ")
		secondAttempt = readPassword("confirm password: ")
		if firstAttempt == secondAttempt {
			break
		}
		fmt.Println("passwords do not match, try again")
	}
	return firstAttempt
}

func readPassword(prompt string) string {
	// TODO: add windows support
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println(err)
	}
	password := string(bytePassword)
	fmt.Println()
	return password
}

func main() {
	var actions = map[string]executor.Action{
		"connect":    {connect, 2},
		"disconnect": {disconnect, 0},
		"start":      {start, 1},
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
