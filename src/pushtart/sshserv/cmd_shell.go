package sshserv

import (
	"io"
	"pushtart/logging"
	"pushtart/sshserv/cmd_registry"
	"pushtart/tartmanager"
	"pushtart/user"
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
			runFunc(util.ParseCommands(util.TokeniseCommandString(line[len(spl[0]):])), &commandOutputRewriter{Out: term})
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
	} else {
		spl := strings.Split(line, " ")

		switch {
		case spl[len(spl)-2] == "--tart":
			t := spl[len(spl)-1]
			if !strings.HasPrefix(t, "/") {
				t = "/" + t
			}
			if match := util.BestPrefixMatch(t, tartmanager.List()); match != "" {
				newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
				return newLine, len(newLine), true
			}

		case spl[len(spl)-2] == "--username":
			if match := util.BestPrefixMatch(spl[len(spl)-1], user.List()); match != "" {
				newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
				return newLine, len(newLine), true
			}
		}

		//now see if we can autocomplete the last value
		if strings.HasPrefix(spl[len(spl)-1], "--") {
			if commandOptions, ok := commandParams[spl[0]]; ok {
				if match := util.BestPrefixMatch(spl[len(spl)-1], commandOptions); match != "" {
					newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
					return newLine, len(newLine), true
				}
			}
		}

	}
	return line, pos, false
}

var commandParams = map[string][]string{
	"edit-user":  []string{"--username", "--password", "--name", "--allow-ssh-password"},
	"make-user":  []string{"--username", "--password", "--name", "--allow-ssh-password"},
	"start-tart": []string{"--tart"},
	"stop-tart":  []string{"--tart"},
	"edit-tart":  []string{"--tart", "--name", "--set-env", "--delete-env", "--log-stdout"},
}

type commandOutputRewriter struct {
	Out io.Writer
}

func (c *commandOutputRewriter) Write(p []byte) (n int, err error) {
	return c.Out.Write([]byte(strings.Replace(string(p), "\n", "\r\n", -1)))
}

var banner = "____            _   _____          _\r\n|  _ \\ _   _ ___| |_|_   _|_ _ _ __| |_\r\n| |_) | | | / __| '_ \\| |/ _` | '__| __|\r\n|  __/| |_| \\__ \\ | | | | (_| | |  | |_\r\n|_|    \\__,_|___/_| |_|_|\\__,_|_|   \\__|\r\nRun 'exit' to exit the shell, and 'help' for a list of commands.\r\n\r\n"
