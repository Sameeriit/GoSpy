package client

import (
	"bufio"
	"net"
	"strings"
)

// ToDo: Move from using a delimiter to message encoding from comms package.
const delim = '\n' // The delimiter used for separating communications.

type GoSpyClient struct {
	conn net.Conn
}

/// readData reads bytes up to receiving a delim from the client.
func (g GoSpyClient) readData() (data string, err error) {
	data, err = bufio.NewReader(g.conn).ReadString(delim)
	if err != nil {
		return "", nil
	}
	data = strings.TrimSuffix(data, "\n")
	return data, nil
}

// sendData sends the given string to the client.
func (g GoSpyClient) sendData(data string) (err error) {
	data += string(delim)
	_, err = g.conn.Write([]byte(data))
	return err
}

// Ping sends a "ping" to the client and waits for a "pong".
func (g GoSpyClient) Ping() (response string, err error) {
	err = g.sendData("ping")
	if err != nil {
		return "", err
	}
	return g.readData()
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
