package sshserv

import (
	"io"
	"pushtart/logging"
	"pushtart/sshserv/cmd_registry"
	"pushtart/util"
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
		logging.Info("sshserv-shell", "["+conn.User()+"]: "+line)
		spl := strings.Split(line, " ")

		if spl[0] == "\\q" || spl[0] == "exit" {
			break
		}

		if ok, runFunc := cmd_registry.Command(spl[0]); ok {
			runFunc(util.ParseCommands(spl[1:]), &commandOutputRewriter{Out: term})
		}
	}
}

func autocomplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if line == "" || key != '\t' {
		return line, pos, false
	}

	//If we are autocompleting command names, find the best match based on command names
	if len(strings.Split(line, " ")) == 1 {
		if match := util.BestPrefixMatch(line, cmd_registry.List()); match != "" {
			return match, len(match), true
		}
	}
	return line, pos, false
}

type commandOutputRewriter struct {
	Out io.Writer
}

func (c *commandOutputRewriter) Write(p []byte) (n int, err error) {
	return c.Out.Write([]byte(strings.Replace(string(p), "\n", "\r\n", -1)))
}

var banner = "____            _   _____          _\r\n|  _ \\ _   _ ___| |_|_   _|_ _ _ __| |_\r\n| |_) | | | / __| '_ \\| |/ _` | '__| __|\r\n|  __/| |_| \\__ \\ | | | | (_| | |  | |_\r\n|_|    \\__,_|___/_| |_|_|\\__,_|_|   \\__|\r\nPress Control-D to exit the command shell.\r\n\r\n"
