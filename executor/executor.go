package executor

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Executable interface {
	AssertExecutable()
}

func Execute(anyType Executable, command string) ([]reflect.Value, error) {
	commandSplit := strings.Split(command, " ")
	commandName := commandSplit[0]
	commandArgs := commandSplit[1:]

	method := reflect.ValueOf(anyType).MethodByName(strings.Title(commandName))
	if !method.IsValid() {
		errorMessage := commandName + " is not a valid action"
		return []reflect.Value{},
			errors.New(errorMessage)
	}
	expectedArgsCnt := method.Type().NumIn()
	givenArgsCnt := len(commandArgs)
	if givenArgsCnt != expectedArgsCnt {
		errorMessage := fmt.Sprintf(
			"wrong number of arguments passed to %s, expected %d, got %d",
			commandName, expectedArgsCnt, givenArgsCnt)
		return []reflect.Value{}, errors.New(errorMessage)
	}

	methodArgs := make([]reflect.Value, givenArgsCnt)
	for idx, _ := range commandArgs {
		methodArgs[idx] = reflect.ValueOf(commandArgs[idx])
	}

	return method.Call(methodArgs), nil
}
