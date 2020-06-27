package commands

import "github.com/psidex/GoSpy/internal/comms"

// ExitSend sends a message to the client for it to stop.
func ExitSend(c comms.Connection) error {
	return c.SendBytes([]byte("exit"))
}

func ExitReply() {}
