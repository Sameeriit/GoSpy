package comms

import (
	"io"
)

// CopyFromConnection takes a Connection and reads bytes from it using RecvBytes, writing them to the io.Writer.
// Stops when any read or write error occurs. Uses the returned channel to pass the error (or nil).
func CopyFromConnection(src Connection, dst io.Writer) <-chan error {
	errChan := make(chan error)
	go func() {
		var err error
		var readBytes []byte
		for {
			readBytes, err = src.RecvBytes()
			if err != nil {
				break
			}
			_, err = dst.Write(readBytes)
			if err != nil {
				break
			}
		}
		errChan <- err
	}()
	return errChan
}

// CopyToConnection takes a io.Reader and reads bytes from it, sending them to the Connection using SendBytes.
// Stops when any read or write error occurs. Uses the returned channel to pass the error (or nil).
func CopyToConnection(src io.Reader, dst Connection) <-chan error {
	errChan := make(chan error)
	go func() {
		_, err := io.Copy(dst, src)
		errChan <- err
	}()
	return errChan
}
