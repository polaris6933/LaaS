package executor

import (
	"fmt"
	"reflect"
	"strings"
)

type Executable interface {
	AssertExecutable()
}

func Execute(anyType Executable, command string) {
	commandSplit := strings.Split(command, " ")
	commandName := commandSplit[0]
	commandArgs := commandSplit[1:]

	method := reflect.ValueOf(anyType).MethodByName(strings.Title(commandName))
	if !method.IsValid() {
		fmt.Println(commandName, "is not a valid action")
		return
	}
	expectedArgsCnt := method.Type().NumIn()
	givenArgsCnt := len(commandArgs)
	if givenArgsCnt != expectedArgsCnt {
		fmt.Printf(
			"wrong number of arguments passed to %s, expected %d, got %d\n",
			commandName, expectedArgsCnt, givenArgsCnt)
		return
	}

	methodArgs := make([]reflect.Value, givenArgsCnt)
	for idx, _ := range commandArgs {
		methodArgs[idx] = reflect.ValueOf(commandArgs[idx])
	}

	method.Call(methodArgs)
}
