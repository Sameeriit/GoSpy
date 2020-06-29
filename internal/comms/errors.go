package comms

import (
	"io"
	"net"
)

// IsConnectionError checks if an error is due to a problem with the connection (a net.Error or io.EOF).
func IsConnectionError(err error) (is bool) {
	if _, is = err.(net.Error); is == true || err == io.EOF {
		return is
	}
	return false
}
