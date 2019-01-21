package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type serverCall struct {
	call   func([]string)
	argCnt int
}

var requests = map[string]serverCall{
	"connect":    {connect, 1},
	"disconnect": {disconnect, 0},
	"start":      {start, 2},
	"kill":       {kill, 1},
	"unlock":     {unlock, 1},
	"save":       {save, 2},
	"pause":      {pause, 1},
	"resume":     {resume, 1},
}

func connect(args []string) {
	fmt.Println("connect()", args)
}

func disconnect(args []string) {
	fmt.Println("disconnect()", args)
}

func start(args []string) {
	fmt.Println("start()", args)
}

func kill(args []string) {
	fmt.Println("kill()", args)
}

func unlock(args []string) {
	fmt.Println("unlock()", args)
}

func save(args []string) {
	fmt.Println("save()", args)
}

func pause(args []string) {
	fmt.Println("pause()", args)
}

func resume(args []string) {
	fmt.Println("resume()", args)
}

func executor(command string) {
	commandSplit := strings.Split(command, " ")
	commandName := commandSplit[0]
	commandArgs := commandSplit[1:]

	request, ok := requests[commandName]
	if !ok {
		fmt.Println(commandName, "is not a valid operation")
		return
	}
	if len(commandArgs) != request.argCnt {
		fmt.Println("wrong number of arguments for", commandName)
		return
	}
	request.call(commandArgs)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	executor(strings.TrimSuffix(input, "\n"))
}
