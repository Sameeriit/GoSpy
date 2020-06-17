package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// ToDo: This is just basic testing code, requires a lot of cleanup.

func main() {
	address := *flag.String("a", "127.0.0.1:12345", "the address (ip:port) of the gospyserver to connect to")
	log.Printf("Got address %s\n", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		message = strings.TrimSuffix(message, "\n")
		log.Printf("Got %s", string(message))

		switch message {
		case "exit":
			log.Println("Bye!")
			os.Exit(0)
		case "ping":
			_, _ = fmt.Fprintf(conn, "pong\n")
		}
	}
}
