package commands

import "github.com/psidex/GoSpy/internal/comms"

// SendPong replies to the "ping" sent from the given ConnectionManager.
func SendPong(cm comms.ConnectionManager) error {
	return cm.SendBytes([]byte("pong"))
}
