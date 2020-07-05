package serverprompt

import (
	"github.com/c-bata/go-prompt"
)

var suggestions = []prompt.Suggest{
	{"ping", "CommandPing the connected client"},
	{"reverse-shell", "Initiate a reverse shell"},
	{"grab-file", "Download a file from the client; grab-file [path on client] [local path]"},
	{"drop-file", "Upload a file to the client; drop-file [local path] [path on client]"},
	{"exit", "Tell the GoSpy client to exit then and exit the GoSpy Server"},
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}
