// Package executor provides a mechanism for running methods of a given type
// providing an instance of this type and the name of the method as a string.
package executor

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// A type must implement this interface in order for the executor to be able
// to work with it. The method is a no-op, it's just needed to mark a type as
// executable.
type Executable interface {
	AssertExecutable()
}

// Return a reflect object which, when executed, will run the method described
// by `anyType` and `command` and nil. `command` should be a string containing
// the method name and the arguments for that method.
// An error is returned if:
//   - `command` is not a method of the type of `anyType`
//   - `command` does not contain enough arguments for the method it describes
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
