package server

import (
	"github.com/psidex/GoSpy/internal/comms"
	"net"
)

// ConMan is a manager for connections to the client.
type ConMan struct {
	listener net.Listener     // The actual listener for connections from the client.
	Conn     comms.Connection // The main connection to the client.
	pwdStr   string           // The password for encrypting data over the connection.
}

// NewConMan instantiates a ConMan, binds a listener to the given address, and waits for a client to connect.
func NewConMan(bindAddress, password string) (s ConMan, err error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return ConMan{}, err
	}
	s = ConMan{listener: l, pwdStr: password}
	s.Conn = s.WaitForConnection()
	return s, nil
}

// WaitForConnection waits for a successful connection to the listener and then sets up and returns a Connection.
func (m ConMan) WaitForConnection() comms.Connection {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			continue
		}
		if m.pwdStr != "" {
			return comms.NewEncryptedConnection(conn, m.pwdStr)
		}
		return comms.NewPlainConnection(conn)
	}
}
