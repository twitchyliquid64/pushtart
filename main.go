package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"pushtart/config"
	"pushtart/dnsserv"
	"pushtart/logging"
	"pushtart/sshserv"
	"pushtart/sshserv/cmd_registry"
	"pushtart/util"
	"pushtart/webproxy"
	"syscall"
)

func help(params map[string]string, w io.Writer, user string) {
	if w == os.Stdout {
		fmt.Fprintln(w, "USAGE: pushtart <command> [--config <config file>] [command-specific-arguments...]")
		fmt.Fprintln(w, "If no config file is specified, config.json will be used.")
		fmt.Fprintln(w, "SSH server keys, user information, tart status, and other (normally external) information is stored in the config file.")
		fmt.Fprintln(w, "Commands:")
		fmt.Fprintln(w, "\trun")
		fmt.Fprintln(w, "\tmake-config")
	}
	fmt.Fprintln(w, "\tget-config-value --field <config-field> (EG: --field DNS.Listener)")
	fmt.Fprintln(w, "\tset-config-value --field <config-field> --value <new-value>")
	fmt.Fprintln(w, " ")
	if w == os.Stdout {
		fmt.Fprintln(w, "\timport-ssh-key --username <username> [--pub-key-file <path-to-.pub-file>] (Not available from SSH shell)")
	}
	fmt.Fprintln(w, "\tmake-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]")
	fmt.Fprintln(w, "\tedit-user --username <username [--password <password] [--name <name] [--allow-ssh-password yes/no]")
	fmt.Fprintln(w, "\tls-users")
	fmt.Fprintln(w, " ")
	fmt.Fprintln(w, "\tls-tarts")
	if w != os.Stdout {
		fmt.Fprintln(w, "\tstart-tart --tart <pushURL>")
		fmt.Fprintln(w, "\tstop-tart --tart <pushURL>")
	}
	fmt.Fprintln(w, "\tedit-tart --tart <pushURL>[--name <name>] [--set-env \"<name>=<value>\"] [--delete-env <name>] [--log-stdout yes/no]")
	fmt.Fprintln(w, "\ttart-add-owner --tart <pushURL> --username <username>")
	fmt.Fprintln(w, "\ttart-remove-owner --tart <pushURL> --username <username>")
	fmt.Fprintln(w, "\textension --extension <extension name> [command-specific-arguments...]")

	if w != os.Stdout {
		fmt.Fprintln(w, "\tlogs")
	}
}

func main() {
	defer config.UnlockConfig() //Only unlocks if a config was successfully locked

	if len(os.Args) < 2 {
		help(nil, os.Stdout, "")
	} else {

		params := parseCommands(os.Args[2:])
		switch os.Args[1] {
		case "run":
			var err error
			logging.Info("main", "Starting in runmode")
			configInit(params["config"])
			registerCommands()
			err = sshserv.Init()
			if err != nil {
				logging.Error("main-sshserv", err.Error())
			}
			dnsserv.Init()
			webproxy.Init()

			c := make(chan os.Signal, 2)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			fmt.Println("")
			logging.Info("main", "Recieved Interrupt, shutting down")

		case "make-config":
			generateConfig(params["config"])

		case "make-user":
			configInit(params["config"])
			makeUser(params, os.Stdout, "")

		case "edit-user":
			configInit(params["config"])
			editUser(params, os.Stdout, "")

		case "ls-users":
			configInit(params["config"])
			listUser(params, os.Stdout, "")

		case "ls-tarts":
			configInit(params["config"])
			listTarts(params, os.Stdout, "")

		case "edit-tart":
			configInit(params["config"])
			editTart(params, os.Stdout, "")

		case "tart-restart-mode":
			configInit(params["config"])
			tartRestartMode(params, os.Stdout, "")

		case "import-ssh-key":
			configInit(params["config"])
			importSSHKey(params, os.Stdout, "")

		case "extension":
			configInit(params["config"])
			extensionCommand(params, os.Stdout, "")

		case "get-config-value":
			configInit(params["config"])
			getConfigValue(params, os.Stdout, "")

		case "set-config-value":
			configInit(params["config"])
			setConfigValue(params, os.Stdout, "")

		case "tart-add-owner":
			configInit(params["config"])
			tartAddOwner(params, os.Stdout, "")

		case "tart-remove-owner":
			configInit(params["config"])
			tartRemoveOwner(params, os.Stdout, "")
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
	cmd_registry.Register("logs", logMsgs)
	cmd_registry.Register("tart-restart-mode", tartRestartMode)
	cmd_registry.Register("extension", extensionCommand)
	cmd_registry.Register("get-config-value", getConfigValue)
	cmd_registry.Register("set-config-value", setConfigValue)
	cmd_registry.Register("tart-add-owner", tartAddOwner)
	cmd_registry.Register("tart-remove-owner", tartRemoveOwner)
}

// configInit loads the configuration file from the command line. If there was an error loading the file, a default configuration
// is generated.
func configInit(configPath string) error {
	var err error
	err = config.Load(configPath)

	if err == config.ErrLockfileExists {
		fmt.Println("Lock exists for another process (" + configPath + ".lock)")
		os.Exit(1)
	}

	exists, _ := util.FileExists(configPath)
	if err != nil && !exists {
		generateConfig(configPath)
	} else if err != nil {
		logging.Fatal("init", "Please review the configuration file for errors, and consider deleting it if you would like an empty one generated on next start.")
	}
	return err
}

func generateConfig(configPath string) {
	err := config.Generate(configPath)

	if err != nil {
		logging.Error("init-generateConfig", err.Error())
	}
}
