package executor

import (
	"fmt"
	"net"
	"strings"
)

type Action struct {
	Call   func([]string, *net.Conn)
	ArgCnt int
}

type Executor struct {
	actions    map[string]Action
	Connection *net.Conn
}

func NewExecutor(actions map[string]Action) *Executor {
	var e Executor
	e.actions = actions
	e.Connection = new(net.Conn)
	return &e
}

func (e *Executor) Execute(action string) {
	actionSplit := strings.Split(action, " ")
	actionName := actionSplit[0]
	actionArgs := actionSplit[1:]

	action, ok := e.actions[actionName]
	if !ok {
		fmt.Println(actionName, "is not a valid action")
		return
	}
	if len(actionArgs) != action.ArgCnt {
		fmt.Println("wrong number of arguments for", actionName)
		return
	}
	action.Call(actionArgs, e.Connection)
}
