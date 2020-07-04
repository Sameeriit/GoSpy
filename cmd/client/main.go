package main

import (
	"flag"
	"github.com/psidex/GoSpy/internal/client"
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"net"
)

func main() {
	address := flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	flag.Parse()

	for {
		log.Printf("Attempting connection to address: %s\n", *address)
		conn, err := net.Dial("tcp", *address)
		if err != nil {
			continue
		}
		// The connection to the server for exchanging command data (similar to CmdCon in server conman).
		cmdCon := comms.NewConnection(conn)
		log.Println("Successful connection")
		err = client.CommandLoop(cmdCon)
		_ = cmdCon.Close()
		log.Printf("Connection dropped: %s\n", err.Error())
	}
}
