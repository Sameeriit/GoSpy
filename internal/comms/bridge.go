package comms

import (
	"io"
)

// BridgeCMToWriter takes a ConnectionManager and a writer, reading bytes from the cm and sending them to the writer.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgeCMToWriter(cm ConnectionManager, dst io.Writer) <-chan error {
	errChannel := make(chan error)
	go func() {
		var err error
		for {
			var readBytes []byte
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

// BridgerReaderToCM takes a reader and reads bytes from it, sending them to the ConnectionManager using SendBytes.
// Uses the returned channel to signify that an error occurred (and to pass it).
func BridgerReaderToCM(src io.Reader, cm ConnectionManager) <-chan error {
	errChannel := make(chan error)
	go func() {
		var err error
		readBuf := make([]byte, 1024)
		for {
			var nBytes int
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
