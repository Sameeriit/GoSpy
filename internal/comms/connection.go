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
	return Connection{conn}
}

// Write implements the io.Writer interface, actually just calls SendBytes.
func (c Connection) Write(buf []byte) (n int, err error) {
	return len(buf), c.SendBytes(buf)
}

// SendBytes sends a slice of bytes over the connection.
func (c Connection) SendBytes(data []byte) (err error) {
	err = binary.Write(c.conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// RecvBytes receives a slice of bytes over the connection that was sent by SendBytes.
func (c Connection) RecvBytes() (data []byte, err error) {
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

// GetRemoteAddr returns the address of the remote connection as a string ("192.0.2.1:25", "[2001:db8::1]:80", etc.).
func (c Connection) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// Close closes the connection.
func (c Connection) Close() error {
	return c.conn.Close()
}
