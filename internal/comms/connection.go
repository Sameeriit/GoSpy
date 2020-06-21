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

func (c PlainConn) SendBytes(data []byte) (err error) {
	return sendBytesTo(c.conn, data)
}

func (c PlainConn) RecvBytes() (data []byte, err error) {
	return recvBytesFrom(c.conn)
}

func (c PlainConn) Close() error {
	return c.conn.Close()
}

// EncryptedConn is a connection manager that encrypts its packets using a provided byteEncryptor.
type EncryptedConn struct {
	conn net.Conn
	be   byteEncryptor
}

// NewEncryptedConn instantiates a new EncryptedConn.
func NewEncryptedConn(conn net.Conn, password string) EncryptedConn {
	be := newByteEncryptor(password)
	return EncryptedConn{conn, be}
}

func (c EncryptedConn) SendBytes(data []byte) (err error) {
	encrypted, err := c.be.Encrypt(data)
	if err != nil {
		return err
	}
	return sendBytesTo(c.conn, encrypted)
}

func (c EncryptedConn) RecvBytes() (data []byte, err error) {
	dataBytes, err := recvBytesFrom(c.conn)
	if err != nil {
		return nil, err
	}

	plaintext, err := c.be.Decrypt(dataBytes)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c EncryptedConn) Close() error {
	return c.conn.Close()
}

// Helper functions.

// sendBytesTo sends data over conn, prepending data's size as a uint64.
func sendBytesTo(conn net.Conn, data []byte) (err error) {
	err = binary.Write(conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

// recvBytesFrom receives a uint64 and then reads that many bytes into a slice and returns it.
func recvBytesFrom(conn net.Conn) (data []byte, err error) {
	var length int64
	err = binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	data = make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
