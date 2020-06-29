package conman

import (
	"github.com/psidex/GoSpy/internal/comms"
	"net"
)

// ConMan is a connection manager for handling connections to and from the client.
type ConMan struct {
	listener net.Listener     // The listener for connections from the client.
	CmdCon   comms.Connection // The connection to the client that exchanges commands and replies.
}

// NewConMan instantiates a ConMan, binds a listener to the given address, and waits for a client to connect.
func NewConMan(bindAddress string) (s ConMan, err error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return ConMan{}, err
	}
	return ConMan{listener: l}, nil
}

// WaitForNewConnection waits for a successful connection to the listener and then sets up and returns a
// comms.Connection.
func (m ConMan) WaitForNewConnection() comms.Connection {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			continue
		}
		return comms.NewConnection(conn)
	}
}

// Stop unbinds the listener and closes the CmdCon (ignores errors).
func (m ConMan) Stop() {
	_ = m.listener.Close()
	_ = m.CmdCon.Close()
}
