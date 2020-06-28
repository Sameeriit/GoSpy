package commands

import (
	"errors"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
)

// PingReply sends "pong" to the given connection.
func PingReply(c comms.Connection) error {
	return c.SendBytes([]byte("pong"))
}

// PingSend sends a ping request to the connection and waits for a pong reply.
func PingSend(c comms.Connection) (err error) {
	if err = c.SendBytes([]byte("ping")); err != nil {
		return err
	}

	var reply []byte
	if reply, err = c.RecvBytes(); err != nil {
		return err
	}

	if string(reply) != "pong" {
		return errors.New("did not receive \"pong\"")
	}

	fmt.Println("Received pong")
	return nil
}
