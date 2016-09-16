package main

import (
	"fmt"
	"io"
	"os"
	"pushtart/config"
	"pushtart/constants"
	"pushtart/logging"
	"pushtart/sshserv"
	"pushtart/sshserv/cmd_registry"
	"time"
)

func help(params map[string]string, w io.Writer) {
	if w == os.Stdout {
		fmt.Fprintln(w, "USAGE: pushtart <command> [--config <config file>] [command-specific-arguments...]")
		fmt.Fprintln(w, "If no config file is specified, config.json will be used.")
		fmt.Fprintln(w, "SSH server keys, user information, tart status, and other (normally external) information is stored in the config file.")
		fmt.Fprintln(w, "Commands:")
		fmt.Fprintln(w, "\trun (Not available from SSH shell)")
		fmt.Fprintln(w, "\tmake-config (Not available from SSH shell)")
		fmt.Fprintln(w, "\timport-ssh-key --username <username> [--pub-key-file <path-to-.pub-file>] (Not available from SSH shell)")
	}
	fmt.Fprintln(w, "\tmake-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]")
	fmt.Fprintln(w, "\tedit-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]")
	fmt.Fprintln(w, "\tls-users")
	fmt.Fprintln(w, "\tls-tarts")
	fmt.Fprintln(w, "\tstart-tart --tart <pushURL>")
	fmt.Fprintln(w, "\tstop-tart --tart <pushURL>")
	fmt.Fprintln(w, "\tedit-tart --tart <pushURL>[--name <name>] [--set-env \"<name>=<value>\"] [--delete-env <name>]")
}

func main() {

	if len(os.Args) < 2 {
		help(nil, os.Stdout)
	} else {

		params := parseCommands(os.Args[2:])
		switch os.Args[1] {
		case "run":
			var err error
			logging.Info("init", "Starting in runmode")
			configInit(params["config"])
			registerCommands()
			err = sshserv.Init()
			if err != nil {
				logging.Error("init-sshServ", err.Error())
			}
			for {
				time.Sleep(1 * time.Second)
			}

		case "make-config":
			generateConfig(params["config"])

		case "make-user":
			configInit(params["config"])
			makeUser(params, os.Stdout)

		case "edit-user":
			configInit(params["config"])
			editUser(params, os.Stdout)

		case "ls-users":
			configInit(params["config"])
			listUser(params, os.Stdout)

		case "ls-tarts":
			configInit(params["config"])
			listTarts(params, os.Stdout)

		case "start-tart":
			configInit(params["config"])
			startTart(params, os.Stdout)

		case "stop-tart":
			configInit(params["config"])
			stopTart(params, os.Stdout)

		case "edit-tart":
			configInit(params["config"])
			editTart(params, os.Stdout)

		case "import-ssh-key":
			configInit(params["config"])
			importSSHKey(params, os.Stdout)
		}
	}
}

func registerCommands() {
	cmd_registry.Register("make-user", makeUser)
	cmd_registry.Register("edit-user", editUser)
	cmd_registry.Register("ls-users", listUser)
	cmd_registry.Register("ls-tarts", listTarts)
	cmd_registry.Register("start-tart", startTart)
	cmd_registry.Register("stop-tart", stopTart)
	cmd_registry.Register("edit-tart", editTart)
	cmd_registry.Register("help", help)
}

// configInit loads the configuration file from the command line. If there was an error loading the file, a default configuration
// is generated.
func configInit(configPath string) error {
	var err error
	if len(os.Args) > 2 {
		err = config.Load(configPath)
	} else {
		err = config.Load(constants.DefaultConfigFileName)
	}
	if err != nil {
		generateConfig(configPath)
	}
	return err
}

func generateConfig(configPath string) {
	var err error
	if len(os.Args) > 2 {
		err = config.Generate(configPath)
	} else {
		err = config.Generate(constants.DefaultConfigFileName)
	}
	if err != nil {
		logging.Error("init-generateConfig", err.Error())
	}
}
