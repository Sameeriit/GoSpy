package main

import (
	"flag"
	"fmt"
	"github.com/psidex/GoSpy/internal/server/conman"
	"github.com/psidex/GoSpy/internal/server/serverprompt"
	"log"
	"os"
)

const banner = `
   ___     ___             ___                      
  / __|___/ __|_ __ _  _  / __| ___ _ ___ _____ _ _ 
 | (_ / _ \__ \ '_ \ || | \__ \/ -_) '_\ V / -_) '_|
  \___\___/___/ .__/\_, | |___/\___|_|  \_/\___|_|  
              |_|   |__/              v0.0.1`

func main() {
	fmt.Printf("%s\n\n", banner)

	bindAddr := flag.String("a", "0.0.0.0:12345", "the address (ip:port) to bind the gospyserver to")
	password := flag.String("p", "", "the password to encrypt network data with")
	flag.Parse()

	if *password != "" {
		fmt.Println("Password supplied, using encrypted connection")
	}

	man, err := conman.NewConMan(*bindAddr, *password)
	if err != nil {
		fmt.Printf("Error binding listener: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Type `exit` or `Ctrl-C` to exit\nListening on %s\n", *bindAddr)
	for {
		fmt.Println("Waiting for connection from GoSpy client...")
		man.CmdCon = man.WaitForNewConnection()
		fmt.Println("Successful connection from client")
		err = serverprompt.CommandLoop(man)
		_ = man.CmdCon.Close()
		log.Printf("Client connection dropped: %s\n", err.Error())
	}
}
