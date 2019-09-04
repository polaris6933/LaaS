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

func (c *Client) makeRequest(requestArgs []string) string {
	if c.connection == nil {
		return "not connected to server atm"
	}
	request := strings.Join(requestArgs, " ")
	request = request + "\000"
	fmt.Println(request)
	fmt.Fprintf(*c.connection, request)

	return c.waitResponse()
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
	// TODO: disallow empty user name
	if username == defaultUserName {
		return "user name " + defaultUserName + " not allowed"
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

func (c *Client) List() string {
	return c.makeRequest([]string{"list"})
}

// TODO: implement actual remove
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
	return "game is not being displayed"
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
	// response = client.Connect("localhost:8001")
	// fmt.Println(response)
	// response = client.Register("ivan")
	// fmt.Println(response)
	// response = client.Login("ivan")
	// fmt.Println(response)
	// response = client.Start("s")
	// fmt.Println(response)
	// response = client.Start("t")
	// fmt.Println(response)
	// response = client.Start("s")
	// fmt.Println(response)
	// response = client.Start("i")
	// fmt.Println(response)
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
