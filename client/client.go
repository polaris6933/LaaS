// This package provides the entry point of the client application for LaaS
// as well as the business logic of the client.
package main

import (
	"LaaS/executor"
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const connectionType = "tcp"
const defaultUserName = "none"
const hostname = "localhost"
const port = ":8088"

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// A client is described by its connection to the server and the name of
// currently logged user. A default username is used to denote no one is
// currently logged in.
// The methods of `Client` return a user readable string describing the result
// of the issued operation. It is either a response from the server or an error
// raised by the client.
type Client struct {
	connection *net.Conn
	loggedAs   string
}

// Constructs a new client.
func NewClient() *Client {
	c := new(Client)
	c.connection = nil
	c.loggedAs = defaultUserName
	return c
}

// Implement the Executable interface.
func (c Client) AssertExecutable() {}

func (c *Client) attemptRecconect() {
	c.loggedAs = defaultUserName
	for {
		fmt.Println("attempting to reconnect")
		connect := c.Connect(hostname + port)
		if connect != "" {
			fmt.Println(connect)
			break
		}
		time.Sleep(time.Second)
	}
}

func (c *Client) makeRequest(requestArgs []string) string {
	if c.connection == nil {
		return "not connected to server atm"
	}
	request := strings.Join(requestArgs, " ")
	request = request + "\000"
	fmt.Fprintf(*c.connection, request)

	response := c.waitResponse()
	if response == "" {
		fmt.Println("connection to the server has been lost")
		c.attemptRecconect()
	}
	return response
}

func (c *Client) waitResponse() string {
	response, err := bufio.NewReader(*c.connection).ReadString('\000')
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.TrimRight(response, "\000")
}

// Makes a request to the server attempting to register a user with the name
// `username`.
// Fails if:
//   - `username` is the default user name
//   - `username` is the empty string
func (c *Client) Register(username string) string {
	if username == defaultUserName {
		return "user name " + defaultUserName + " not allowed"
	} else if username == "" {
		return "uesr name must not be empty"
	}
	password := passwordConfirmation()
	response := c.makeRequest([]string{"register", username, password})
	if response == "registered user "+username {
		c.loggedAs = username
	}
	return response
}

// Makes a request to the server attempting to log the user in.
func (c *Client) Login(username string) string {
	password := readPassword("input password: ")
	response := c.makeRequest([]string{"login", username, password})
	if response == "user "+username+" logged in" {
		c.loggedAs = username
	}
	return response
}

// Logs the user out.
func (c *Client) Logout() string {
	c.loggedAs = defaultUserName
	return "logged out"
}

// Attempts to establish a connection to the server.
func (c *Client) Connect(connectTo string) string {
	c.Disconnect()
	conn, err := net.Dial(connectionType, connectTo)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	c.connection = &conn
	return "connected to " + (*c.connection).RemoteAddr().String()
}

// Disconnects from the server.
func (c *Client) Disconnect() string {
	if c.connection == nil {
		return ""
	}
	address := (*c.connection).RemoteAddr().String()
	(*c.connection).Close()
	c.connection = nil
	return "disconnected from " + address
}

// Makes a request to the server attempting to add a new session.
func (c *Client) Add(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"add", c.loggedAs, name})
}

// Makes a request to the server attempting to start a session with the given
// configuration.
// Fails if the user is not logged in.
func (c *Client) Start(name, config string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"start", c.loggedAs, name, config})
}

// Makes a request to the server attempting to resume a stopped session.
// Fails if the user is not logged in.
func (c *Client) Resume(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"resume", c.loggedAs, name})
}

// Makes a request to the server attempting to list the session on the server.
func (c *Client) List() string {
	return c.makeRequest([]string{"list"})
}

// Makes a request to the server attempting to stop a session.
// Fails if the user is not logged in.
func (c *Client) Stop(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"stop", c.loggedAs, name})
}

// Makes a request to the server attempting to delete a session.
// Fails if the user is not logged in.
func (c *Client) Kill(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"kill", c.loggedAs, name})
}

func (c *Client) displayGame(s chan os.Signal, name string) {
	for {
		select {
		case <-s:
			clearScreen()
			signal.Reset(os.Interrupt)
			return
		default:
			board := c.makeRequest([]string{"watch", c.loggedAs, name})
			clearScreen()
			fmt.Println(board)
			time.Sleep(time.Second)
		}
	}
}

// Makes a request to the server attempting to retrieve the current state of
// the game associated with a session.
// Fails if the user is not logged in.
func (c *Client) Watch(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	s := make(chan os.Signal, 2)
	signal.Notify(s, os.Interrupt)
	go c.displayGame(s, name)
	return "game is now being displayed"
}

// Stops the application.
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
	fmt.Print(prompt)
	var bytePassword []byte
	var err error
	for {
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println(err)
			return ""
		}
		if len(bytePassword) > 0 {
			break
		}
	}
	fmt.Println()
	password := string(bytePassword)
	return password
}

func main() {
	client := NewClient()
	var response string
	defer client.Disconnect()
	connect := client.Connect(hostname + port)
	fmt.Println(connect)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		execResult, err := executor.Execute(
			client, strings.TrimSpace(strings.TrimSuffix(input, "\n")))
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			response = execResult[0].Interface().(string)
		}
		fmt.Println(response)
	}
}
