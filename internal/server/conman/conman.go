package conman

import (
	"github.com/psidex/GoSpy/internal/comms"
	"net"
)

// ConMan is a connection manager for handling connections to/from the client.
type ConMan struct {
	listener net.Listener     // The listener for connections from the client.
	CmdCon   comms.Connection // The connection to the client for exchanging command data (e.g. "do this command").
}

// NewConMan instantiates a ConMan and binds a listener to the given address.
func NewConMan(bindAddress string) (newCM ConMan, err error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return ConMan{}, err
	}
	return ConMan{listener: l}, nil
}

// AcceptCmdCon calls m.AcceptSuccessful and then assigns the returned net.Conn to the m.CmdCon field (as a Connection).
func (m *ConMan) AcceptCmdCon() {
	newConn := m.AcceptSuccessful()
	m.CmdCon = comms.NewConnection(newConn)
}

// AcceptSuccessful waits for a successful connection to the listener and returns the net.Conn.
func (m ConMan) AcceptSuccessful() net.Conn {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			continue
		}
		return conn
	}
}

// Stop unbinds the listener and closes the CmdCon (ignores errors).
func (m ConMan) Stop() {
	_ = m.listener.Close()
	_ = m.CmdCon.Close()
}
