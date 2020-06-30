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

// NewConnectionToRemote creates a new Connection to the remote address of the current.
func (c Connection) NewConnectionToRemote() (Connection, error) {
	conn, err := net.Dial("tcp", c.conn.RemoteAddr().String())
	if err != nil {
		return Connection{}, err
	}
	return NewConnection(conn), nil
}

// Write implements the io.Writer interface, actually just calls SendBytes.
func (c Connection) Write(buf []byte) (n int, err error) {
	return len(buf), c.SendBytes(buf)
}

// Close closes the connection.
func (c Connection) Close() error {
	return c.conn.Close()
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

// GoReadFrom takes a io.Reader and reads bytes from it, sending them to the Connection using SendBytes. It does this
// copying in a goroutine and uses the returned channel to pass an error if it occurs. The channel receiving an error
// (or nil) signifies the end of the goroutine.
func (c Connection) GoReadFrom(src io.Reader) <-chan error {
	errChan := make(chan error)
	go func() {
		_, err := io.Copy(c, src)
		errChan <- err
	}()
	return errChan
}

// GoWriteTo takes a io.Reader and reads bytes from it, sending them to the Connection using SendBytes. It does this
// copying in a goroutine and uses the returned channel to pass an error if it occurs. The channel receiving an error
// (or nil) signifies the end of the goroutine.
func (c Connection) GoWriteTo(dst io.Writer) <-chan error {
	errChan := make(chan error)
	go func() {
		var err error
		var readBytes []byte
		for {
			if readBytes, err = c.RecvBytes(); err != nil {
				break
			}
			if _, err = dst.Write(readBytes); err != nil {
				break
			}
		}
		errChan <- err
	}()
	return errChan
}
