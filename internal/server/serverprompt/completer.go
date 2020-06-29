package serverprompt

import (
	"github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	{"ping", "CommandPing the connected client"},
	{"reverse-shell", "Initiate a reverse shell"},
	{"grab-file", "Download a file from the client, grab-file [path on client] [local path]"},
	{"exit", "Tell the GoSpy client to exit then and exit the GoSpy Server"},
}

// Completer is the completer function for the go-prompt prompt.
func Completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}
