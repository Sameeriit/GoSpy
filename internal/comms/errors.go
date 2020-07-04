package comms

import (
	"io"
	"net"
)

// IsNetworkError checks if an error is due to a problem with the connection (a net.Error or io.EOF).
func IsNetworkError(err error) bool {
	if _, isNetErr := err.(net.Error); isNetErr == true || err == io.EOF {
		return true
	}
	return false
}
