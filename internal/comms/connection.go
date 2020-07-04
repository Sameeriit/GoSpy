package comms

import (
	"encoding/binary"
	"io"
	"net"
)

// Connection is for sending and receiving "packets" of bytes over a net.Conn TCP connection.
type Connection struct {
	conn net.Conn
}

// NewConnection instantiates a new Connection.
func NewConnection(conn net.Conn) Connection {
	// Use a "New" function allows us to keep conn unexported.
	return Connection{conn}
}

// Close closes the connection.
func (c Connection) Close() error {
	return c.conn.Close()
}

// DialRemote creates a new net.Conn connection to the remote address of the current Connection.
func (c Connection) DialRemote() (conn net.Conn, err error) {
	return net.Dial("tcp", c.conn.RemoteAddr().String())
}

// sendBytes sends a slice of bytes over the connection.
func (c Connection) sendBytes(data []byte) (err error) {
	err = binary.Write(c.conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// recvBytes receives a slice of bytes over the connection that was sent by sendBytes.
func (c Connection) recvBytes() (data []byte, err error) {
	var length int64
	err = binary.Read(c.conn, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	data = make([]byte, length)
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SendString sends a string over the connection.
func (c Connection) SendString(str string) error {
	return c.sendBytes([]byte(str))
}

// RecvString receives a string the connection that was sent by SendString.
func (c Connection) RecvString() (string, error) {
	bytes, err := c.recvBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
