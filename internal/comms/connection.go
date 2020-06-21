package comms

import (
	"encoding/binary"
	"io"
	"net"
)

// ConnectionManager defines an interface for sending and receiving "packets" of bytes over a net.Conn TCP connection.
type ConnectionManager interface {
	SendBytes([]byte) error
	RecvBytes() ([]byte, error)
	Close() error
}

// PlainConn is a plaintext connection manager.
type PlainConn struct {
	conn net.Conn
}

// NewPlainConn instantiates a new PlainConn.
func NewPlainConn(conn net.Conn) PlainConn {
	return PlainConn{conn}
}

// SendBytes sends a slice of bytes over the connection.
func (c PlainConn) SendBytes(data []byte) (err error) {
	err = binary.Write(c.conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// RecvBytes receives a slice of bytes over the connection that was sent by SendBytes.
func (c PlainConn) RecvBytes() (data []byte, err error) {
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

// Close closes the connection.
func (c PlainConn) Close() error {
	return c.conn.Close()
}

// EncryptedConn is a connection manager that encrypts its packets using a provided byteEncryptor.
type EncryptedConn struct {
	PlainConn
	be byteEncryptor
}

// NewEncryptedConn instantiates a new EncryptedConn.
func NewEncryptedConn(conn net.Conn, password string) EncryptedConn {
	be := newByteEncryptor(password)
	return EncryptedConn{PlainConn{conn}, be}
}

// SendBytes encrypts data and then sends its using PlainConn.SendBytes.
func (c EncryptedConn) SendBytes(data []byte) (err error) {
	encrypted, err := c.be.Encrypt(data)
	if err != nil {
		return err
	}
	return c.PlainConn.SendBytes(encrypted)
}

// RecvBytes receives bytes using PlainConn.RecvBytes and then decrypts them.
func (c EncryptedConn) RecvBytes() (data []byte, err error) {
	encryptedBytes, err := c.PlainConn.RecvBytes()
	if err != nil {
		return nil, err
	}

	data, err = c.be.Decrypt(encryptedBytes)
	if err != nil {
		return nil, err
	}

	return data, nil
}
