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
			runFunc(util.ParseCommands(util.TokeniseCommandString(line[len(spl[0]):])), &commandOutputRewriter{Out: term}, conn.User())
		}
	}
}

type commandOutputRewriter struct {
	Out io.Writer
}

func (c *commandOutputRewriter) Write(p []byte) (n int, err error) {
	return c.Out.Write([]byte(strings.Replace(string(p), "\n", "\r\n", -1)))
}

var banner = "____            _   _____          _\r\n|  _ \\ _   _ ___| |_|_   _|_ _ _ __| |_\r\n| |_) | | | / __| '_ \\| |/ _` | '__| __|\r\n|  __/| |_| \\__ \\ | | | | (_| | |  | |_\r\n|_|    \\__,_|___/_| |_|_|\\__,_|_|   \\__|\r\nRun 'exit' to exit the shell, and 'help' for a list of commands.\r\n\r\n"
