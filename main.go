package main

import (
	"fmt"
	"os"
	"strings"
	"pushtart/config"
	"pushtart/constants"
	"pushtart/logging"
	"pushtart/sshserv"
	"pushtart/user"
	"pushtart/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: pushtart <command> [command-specific-arguments...]")
		fmt.Println("If no config file is specified, config.json will be used.")
		fmt.Println("SSH server keys, user information, and other (normally external) information is stored in the config file.")
		fmt.Println("Commands:")
		fmt.Println("\trun [--config <path-to-configuration-file>]")
		fmt.Println("\tmake-config [--config <config file> --parameter-name parameter-value ...]")
		fmt.Println("\tmake-user --username <username [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
		fmt.Println("\tedit-user --username <username [--config <config file>] [--password <password] [--name <name] [--allow-ssh-password yes/no]")
	} else {

		params := parseCommands(os.Args[2:])
		switch os.Args[1] {
		case "run":
			var err error
			logging.Info("init", "Starting in runmode")
			configInit(params["config"])
			err = sshserv.Init()
			if err != nil {
				logging.Error("init-sshServ", err.Error())
			}

		case "make-config":
			generateConfig(params["config"])

		case "make-user":
			configInit(params["config"])
			makeUser(params)

		case "edit-user":
			configInit(params["config"])
			editUser(params)
		}
	}
}



func saveUser(username string, usr config.User, params map[string]string){
	var exists bool

	if _, exists = params["password"]; exists {
		pw, err := util.HashPassword(username, params["password"])
		if err != nil {
			logging.Error("make-user", "Error hashing password: " + err.Error())
			return
		}
		usr.Password = pw
	}

	if _, exists = params["name"]; exists {
		usr.Name = params["name"]
	}

	if _, exists = params["allow-ssh-password"]; exists {
		usr.AllowSSHPassword = false
		if strings.ToUpper(params["allow-ssh-password"]) == "YES" {
			usr.AllowSSHPassword = true
		}
	}

	user.Save(username, usr)
}


func editUser(params map[string]string){
	if missingFields := checkHasFields([]string{"username"}, params); len(missingFields) > 0 {
		fmt.Println("USAGE: pushtart edit-user --username <username>")
		printMissingFields(missingFields)
		return
	}

	usr := user.Get(params["username"])
	saveUser(params["username"], usr, params)
}

func makeUser(params map[string]string){
	if missingFields := checkHasFields([]string{"username"}, params); len(missingFields) > 0 {
		fmt.Println("USAGE: pushtart make-user --username <username>")
		printMissingFields(missingFields)
		return
	}

	if user.Exists(params["username"]){
		fmt.Println("Err: user already exists")
		return
	}

	user.New(params["username"])
	usr := user.Get(params["username"])
	saveUser(params["username"], usr, params)
}

// configInit loads the configuration file from the command line. If there was an error loading the file, a default configuration
// is generated.
func configInit(configPath string)error{
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
