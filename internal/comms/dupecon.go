package comms

import (
	"net"
)

// DupeCon takes a connection (c) and creates a new connection to the remote address of c, re-using the password if it
// was encrypted.
func DupeCon(c Connection) (newCon Connection, err error) {
	conn, err := net.Dial("tcp", c.GetRemoteAddr())
	if err != nil {
		return nil, err
	}

	if ec, ok := c.(EncryptedConnection); ok == true {
		newCon = NewEncryptedConnection(conn, ec.GetPassword())
	} else {
		newCon = NewPlainConnection(conn)
	}

	return newCon, nil
}
