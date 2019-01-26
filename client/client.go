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

type Client struct {
	connection *net.Conn
}

func NewClient() *Client {
	c := new(Client)
	c.connection = nil
	return c
}

func (c Client) AssertExecutable() {}

func (c *Client) WaitResponse() string {
	response, err := bufio.NewReader(*c.connection).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimSuffix(response, "\n")
}

func (c *Client) Connect(connectionType, connectTo string) string {
	c.Disconnect()
	conn, err := net.Dial(connectionType, connectTo)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	c.connection = &conn
	return "connected to " + (*c.connection).RemoteAddr().String()
}

func (c *Client) Disconnect() string {
	if c.connection == nil {
		return ""
	}
	address := (*c.connection).RemoteAddr().String()
	(*c.connection).Close()
	c.connection = nil
	return "disconnected from " + address
}

func (c *Client) Start(name string) string {
	if c.connection == nil {
		return "not connected to server atm"
	}

	password := passwordConfirmation()
	fmt.Fprintf(*c.connection, "add "+name+" "+password+"\n")
	return c.WaitResponse()
}

func (c *Client) List() {
	if c.connection == nil {
		fmt.Println("not connected to server atm")
		return
	}
	fmt.Fprintf(*c.connection, "list")
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

func (c *Client) Exit() {
	c.Disconnect()
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
	client := NewClient()
	var response string
	defer client.Disconnect()
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		execResult, err := executor.Execute(
			client, strings.TrimSuffix(input, "\n"))
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			response = execResult[0].Interface().(string)
		}
		fmt.Println(response)
	}
}
