package comms

import (
	"encoding/binary"
	"io"
	"net"
)

// Connection defines an interface for sending and receiving "packets" of bytes over a net.Conn TCP connection.
type Connection interface {
	SendBytes([]byte) error
	RecvBytes() ([]byte, error)
	Close() error
}

// PlainConnection is a plaintext connection.
type PlainConnection struct {
	conn net.Conn
}

// NewPlainConnection instantiates a new PlainConnection.
func NewPlainConnection(conn net.Conn) PlainConnection {
	return PlainConnection{conn}
}

// SendBytes sends a slice of bytes over the connection.
func (c PlainConnection) SendBytes(data []byte) (err error) {
	err = binary.Write(c.conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

// RecvBytes receives a slice of bytes over the connection that was sent by SendBytes.
func (c PlainConnection) RecvBytes() (data []byte, err error) {
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
func (c PlainConnection) Close() error {
	return c.conn.Close()
}

// EncryptedConnection is a connection that encrypts its packets using a byteEncryptor.
type EncryptedConnection struct {
	PlainConnection
	be byteEncryptor
}

// NewEncryptedConnection instantiates a new EncryptedConnection.
func NewEncryptedConnection(conn net.Conn, password string) EncryptedConnection {
	be := newByteEncryptor(password)
	return EncryptedConnection{PlainConnection{conn}, be}
}

// SendBytes encrypts data and then sends its using PlainConnection.SendBytes.
func (c EncryptedConnection) SendBytes(data []byte) (err error) {
	encrypted, err := c.be.Encrypt(data)
	if err != nil {
		return err
	}
	return c.PlainConnection.SendBytes(encrypted)
}

// RecvBytes receives bytes using PlainConnection.RecvBytes and then decrypts them.
func (c EncryptedConnection) RecvBytes() (data []byte, err error) {
	encryptedBytes, err := c.PlainConnection.RecvBytes()
	if err != nil {
		return nil, err
	}

	data, err = c.be.Decrypt(encryptedBytes)
	if err != nil {
		return nil, err
	}

	return data, nil
}
