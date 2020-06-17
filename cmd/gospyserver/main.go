package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/psidex/GoSpy/internal/gospyserver/client"
	"log"
	"os"
	"strings"
)

const version = "0.0.1"

var (
	goSpyClient client.GoSpyClient
	suggestions = []prompt.Suggest{
		{"ping", "Ping the connected client"},
		{"exit", "Exit GoSpyServer"},
	}
)

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "ping":
		resp, err := goSpyClient.Ping()
		if err != nil {
			fmt.Printf("Got error: %e", err)
		}
		fmt.Printf("Got %s", resp)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func main() {
	address := *flag.String("b", "0.0.0.0:12345", "the address (ip:port) to bind the gospyserver to")

	fmt.Printf("GoSpyServer v%s\n", version)
	fmt.Printf("Listening on %s\n", address)
	fmt.Println("Waiting for connection from GoSpy client...")

	var err error // So we don't assign a local goSpyClient using := below.
	goSpyClient, err = client.GetGoSpyClient(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successful connection")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("gospyserver-repl"),
	)
	p.Run()
}
