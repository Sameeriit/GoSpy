package comms

import (
	"io"
	"net"
	"os"
)

// TransferFile transfers a file over the conn. If send is true, open the path for reading and send the read file over
// conn, otherwise open the path for writing and receives the file contents from conn. To be used in conjunction with
// another call to the same function (e.g. client calls with send=true and server calls with send=false).
func TransferFile(conn net.Conn, localFilePath string, send bool) (err error) {
	var fd *os.File
	if send {
		fd, err = os.Open(localFilePath)
	} else {
		fd, err = os.Create(localFilePath)
	}
	if err != nil {
		return err
	}
	defer fd.Close()

	var writer io.Writer
	var reader io.Reader
	if send {
		reader = fd
		writer = conn
	} else {
		reader = conn
		writer = fd
	}

	_, err = io.Copy(writer, reader)
	return err
}
