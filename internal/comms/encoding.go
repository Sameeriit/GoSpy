package comms

import (
	"encoding/binary"
	"io"
	"net"
)

// SendStringTo takes a conn and a string and sends that string over the conn, prepending a uint64 of the size of the string.
func SendStringTo(conn net.Conn, data string) (err error) {
	err = binary.Write(conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(data))
	return err
}

// RecvStringFrom takes a conn and using the encoding described in SendStringTo, receives a message and casts it to a string.
func RecvStringFrom(conn net.Conn) (data string, err error) {
	var length int64
	err = binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}

	dataBytes := make([]byte, length)
	_, err = io.ReadFull(conn, dataBytes)
	if err != nil {
		return "", err
	}

	return string(dataBytes), nil
}
