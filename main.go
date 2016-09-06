package main

import (
	"fmt"
	"os"
	"pushtart/config"
	"pushtart/constants"
	"pushtart/logging"
	"pushtart/sshserv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: pushtart <command> [command-specific-arguments...]")
		fmt.Println("Commands:")
		fmt.Println("\trun [<path-to-configuration-file>]")
		fmt.Println("\tmake-configuration parameter-name parameter-value ...")
	} else {

		switch os.Args[1] {
		case "run":
			var err error
			logging.Info("init", "Now starting for run")
			if len(os.Args) > 2 {
				err = config.Load(os.Args[2])
			} else {
				err = config.Load(constants.DefaultConfigFileName)
			}
			if err != nil {
				generateConfig()
			}

			err = sshserv.Init()
			if err != nil {
				logging.Error("init-sshServ", err.Error())
			}

		case "make-configuration":
			generateConfig()
		}
	}
}

func generateConfig() {
	var err error
	if len(os.Args) > 2 {
		err = config.Generate(os.Args[2])
	} else {
		err = config.Generate(constants.DefaultConfigFileName)
	}
	if err != nil {
		logging.Error("init-generateConfig", err.Error())
	}
}
