package sshserv

import (
	"pushtart/logging"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func shell(conn *ssh.ServerConn, channel ssh.Channel) {
	term := terminal.NewTerminal(channel, "> ")
	term.Write([]byte(banner))

	term.AutoCompleteCallback = autocomplete

	defer channel.Close()
	for {
		line, err := term.ReadLine()
		if err != nil {
			logging.Error("sshserv-shell", err.Error())
			break
		}
		logging.Info("sshserv-shell", "Got line: "+line)
	}
}

func autocomplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if line == "" || key != '\t' {
		return line, pos, false
	}

	if len(strings.Split(line, " ")) == 1 { //TODO: Refactor this routine to a util function.
		//iterate all the commands, if there is one suffix match then use it.
		matches := 0
		lastMatch := ""
		for _, command := range availableCommands {
			if strings.HasPrefix(command, line) {
				matches++
				lastMatch = command
			}
		}

		if matches == 1 {
			return lastMatch, len(lastMatch), true
		}
	}
	return line, pos, false
}

var availableCommands = []string{"make-user", "list", "edit-user"}

var banner = "____            _   _____          _\r\n|  _ \\ _   _ ___| |_|_   _|_ _ _ __| |_\r\n| |_) | | | / __| '_ \\| |/ _` | '__| __|\r\n|  __/| |_| \\__ \\ | | | | (_| | |  | |_\r\n|_|    \\__,_|___/_| |_|_|\\__,_|_|   \\__|\r\nPress Control-D to exit the command shell.\r\n\r\n"
