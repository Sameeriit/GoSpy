package comms

import (
	"io"
)

// BridgeConnectionManagerToWriter takes a ConnectionManager and reads bytes from it, sending them to the io.Writer.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgeConnectionManagerToWriter(cm ConnectionManager, dst io.Writer) <-chan error {
	errChannel := make(chan error)
	go func() {
		var err error
		var readBytes []byte
		for {
			readBytes, err = cm.RecvBytes()
			if err != nil {
				break
			}
			_, err = dst.Write(readBytes)
			if err != nil {
				break
			}
		}
		_ = cm.Close()
		errChannel <- err
	}()
	return errChannel
}

// BridgeReaderToConnectionManager takes a io.Reader and reads bytes from it, sending them using the ConnectionManager.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgeReaderToConnectionManager(src io.Reader, cm ConnectionManager) <-chan error {
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
			err = cm.SendBytes(readBuf[0:nBytes])
			if err != nil {
				break
			}
		}
		_ = cm.Close()
		errChannel <- err
	}()
	return errChannel
}
