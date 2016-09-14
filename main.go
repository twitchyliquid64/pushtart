package main

import (
	"fmt"
	"os"
	"pushtart/config"
	"pushtart/constants"
	"pushtart/logging"
	"pushtart/sshserv"
	"pushtart/sshserv/cmd_registry"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("USAGE: pushtart <command> [command-specific-arguments...]")
		fmt.Println("If no config file is specified, config.json will be used.")
		fmt.Println("SSH server keys, user information, tart status, and other (normally external) information is stored in the config file.")
		fmt.Println("Commands:")
		fmt.Println("\trun [--config <path-to-configuration-file>]")
		fmt.Println("\tmake-config [--config <config file> --parameter-name parameter-value ...]")
		fmt.Println("\tmake-user --username <username [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
		fmt.Println("\tedit-user --username <username [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
		fmt.Println("\tls-users [--config <config file>]")
		fmt.Println("\tls-tarts [--config <config file>]")
		fmt.Println("\tstart-tart --tart <pushURL> [--config <config file>]")
		fmt.Println("\tstop-tart --tart <pushURL> [--config <config file>]")
		fmt.Println("\timport-ssh-key --username <username> [--pub-key-file <path-to-.pub-file>]")
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
