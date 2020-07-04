package commands

import "github.com/psidex/GoSpy/internal/comms"

// ExitSend sends a message to the client for it to stop.
func ExitSend(cmdCon comms.Connection) error {
	return cmdCon.SendString("exit")
}

// No need for ExitReply as the client just quits.
