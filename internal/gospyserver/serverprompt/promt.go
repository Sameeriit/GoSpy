package serverprompt

import "github.com/c-bata/go-prompt"

var suggestions = []prompt.Suggest{
	{"ping", "Ping the connected client"},
	{"reverse-shell", "Initiate a reverse shell"},
	{"exit", "Exit GoSpyServer"},
}

// Completer is the completer function for the go-prompt prompt used by GoSpyServer.
func Completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}
