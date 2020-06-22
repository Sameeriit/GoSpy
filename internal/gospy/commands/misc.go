package commands

import "github.com/psidex/GoSpy/internal/comms"

// SendPong replies to the "ping" sent from the given Connection.
func SendPong(c comms.Connection) error {
	return c.SendBytes([]byte("pong"))
}
