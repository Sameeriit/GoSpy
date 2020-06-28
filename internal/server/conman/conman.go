package conman

import (
	"github.com/psidex/GoSpy/internal/comms"
	"net"
)

// ConMan is a connection manager for handling connections to and from the client.
type ConMan struct {
	listener net.Listener     // The listener for connections from the client.
	CmdCon   comms.Connection // The connection to the client that exchanges commands and replies.
	pwdStr   string           // The password for encrypting data over all connections.
}

// NewConMan instantiates a ConMan, binds a listener to the given address, and waits for a client to connect.
func NewConMan(bindAddress, password string) (s ConMan, err error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return ConMan{}, err
	}
	s = ConMan{listener: l, pwdStr: password}
	s.CmdCon = s.WaitForNewConnection()
	return s, nil
}

// WaitForNewConnection waits for a successful connection to the listener and then sets up and returns a comms.Connection.
func (m ConMan) WaitForNewConnection() comms.Connection {
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
