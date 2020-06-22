package comms

import (
	"io"
)

// BridgeConnectionToWriter takes a Connection and reads bytes from it, sending them to the io.Writer.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgeConnectionToWriter(c Connection, dst io.Writer) <-chan error {
	errChannel := make(chan error)
	go func() {
		var err error
		var readBytes []byte
		for {
			readBytes, err = c.RecvBytes()
			if err != nil {
				break
			}
			_, err = dst.Write(readBytes)
			if err != nil {
				break
			}
		}
		_ = c.Close()
		errChannel <- err
	}()
	return errChannel
}

// BridgeReaderToConnection takes a io.Reader and reads bytes from it, sending them using the Connection.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgeReaderToConnection(src io.Reader, c Connection) <-chan error {
	errChannel := make(chan error)
	go func() {
		var err error
		var nBytes int
		readBuf := make([]byte, 1024)
		for {
			nBytes, err = src.Read(readBuf)
			if err != nil {
				break
			}
			err = c.SendBytes(readBuf[0:nBytes])
			if err != nil {
				break
			}
		}
		_ = c.Close()
		errChannel <- err
	}()
	return errChannel
}
