package commands

import (
	"errors"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
)

// PingReply sends "pong" to the given connection.
func PingReply(cmdCon comms.Connection) error {
	return cmdCon.SendString("pong")
}

// PingSend sends a ping request to the connection and waits for a pong reply.
func PingSend(cmdCon comms.Connection) (err error) {
	if err = cmdCon.SendString("ping"); err != nil {
		return err
	}

	if reply, err := cmdCon.RecvString(); err != nil {
		return err
	} else if reply != "pong" {
		return errors.New("did not receive \"pong\"")
	}

	fmt.Println("Received pong")
	return nil
}
