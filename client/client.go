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

type Client struct {
	connection *net.Conn
	loggedAs   string
}

func NewClient() *Client {
	c := new(Client)
	c.connection = nil
	c.loggedAs = defaultUserName
	return c
}

func (c Client) AssertExecutable() {}

func (c *Client) attemptRecconect() {
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

func (c *Client) Login(username string) string {
	password := readPassword("input password: ")
	response := c.makeRequest([]string{"login", username, password})
	if response == "user "+username+" logged in" {
		c.loggedAs = username
	}
	return response
}

func (c *Client) Logout() string {
	c.loggedAs = defaultUserName
	return "logged out"
}

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

func (c *Client) Disconnect() string {
	if c.connection == nil {
		return ""
	}
	address := (*c.connection).RemoteAddr().String()
	(*c.connection).Close()
	c.connection = nil
	return "disconnected from " + address
}

func (c *Client) Add(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"add", c.loggedAs, name})
}

func (c *Client) Start(name, config string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"start", c.loggedAs, name, config})
}

func (c *Client) Resume(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"resume", c.loggedAs, name})
}

func (c *Client) List() string {
	return c.makeRequest([]string{"list"})
}

func (c *Client) Stop(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	return c.makeRequest([]string{"stop", c.loggedAs, name})
}

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

func (c *Client) Watch(name string) string {
	if c.loggedAs == defaultUserName {
		return "not logged in"
	}
	s := make(chan os.Signal, 2)
	signal.Notify(s, os.Interrupt)
	go c.displayGame(s, name)
	return "game is now being displayed"
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
