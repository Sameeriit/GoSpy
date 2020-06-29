package comms

import "net"

// DupeCon takes a connection (c) and creates and returns a new connection to the remote address of c.
func DupeCon(c Connection) (newCon Connection, err error) {
	conn, err := net.Dial("tcp", c.GetRemoteAddr())
	if err != nil {
		return Connection{}, err
	}
	return NewConnection(conn), nil
}
