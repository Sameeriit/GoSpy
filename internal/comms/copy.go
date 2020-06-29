package comms

import (
	"io"
)

// CopyFromConnection takes a Connection and reads bytes from it using RecvBytes, writing them to the io.Writer.
// It does this copying in a goroutine and uses the returned channel to pass an error when it occurs.
// The channel receiving an error (or nil) signifies the end of the goroutine.
func CopyFromConnection(src Connection, dst io.Writer) <-chan error {
	errChan := make(chan error)
	go func() {
		var err error
		var readBytes []byte
		for {
			if readBytes, err = src.RecvBytes(); err != nil {
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

// CopyToConnection takes a io.Reader and reads bytes from it, sending them to the Connection using SendBytes.
// It does this copying in a goroutine and uses the returned channel to pass an error if it occurs.
// The channel receiving an error (or nil) signifies the end of the goroutine.
func CopyToConnection(src io.Reader, dst Connection) <-chan error {
	errChan := make(chan error)
	go func() {
		_, err := io.Copy(dst, src)
		errChan <- err
	}()
	return errChan
}
