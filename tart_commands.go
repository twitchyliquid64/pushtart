package main

import (
	"fmt"
	"io"
	"pushtart/config"
	"pushtart/tartmanager"
	"strconv"
	"strings"
)

func listTarts(params map[string]string, w io.Writer) {
	for pushURL, tart := range config.All().Tarts {
		fmt.Fprint(w, tart.Name+" ("+pushURL+"): ")
		if tart.IsRunning {
			fmt.Fprintln(w, "Running (PID "+strconv.Itoa(tart.PID)+")")
		} else {
			fmt.Fprintln(w, "Stopped.")
		}
		if len(tart.Env) > 0 {
			for _, env := range tart.Env {
				fmt.Fprintln(w, "\t"+env)
			}
		}
	}
}

func startTart(params map[string]string, w io.Writer) {
	if missingFields := checkHasFields([]string{"tart"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart start-tart --tart <pushURL>")
		printMissingFields(missingFields, w)
		return
	}

	exists, tart := findTart(params["tart"])
	if !exists {
		fmt.Fprintln(w, "Err: A tart by that pushURL does not exist")
		return
	}
	err := tartmanager.Start(tart.PushURL)
	if err != nil {
		fmt.Fprintln(w, "Err:", err)
	}
}

func stopTart(params map[string]string, w io.Writer) {
	if missingFields := checkHasFields([]string{"tart"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart start-tart --tart <pushURL>")
		printMissingFields(missingFields, w)
		return
	}

	exists, tart := findTart(params["tart"])
	if !exists {
		fmt.Fprintln(w, "Err: A tart by that pushURL does not exist")
		return
	}
	err := tartmanager.Stop(tart.PushURL)
	if err != nil {
		fmt.Fprintln(w, "Err:", err)
	}
}

func findTart(tartName string) (bool, config.Tart) {
	if tartmanager.Exists(tartName) {
		return true, tartmanager.Get(tartName)
	}
	if tartmanager.Exists("/" + tartName) {
		return true, tartmanager.Get("/" + tartName)
	}
	return false, config.Tart{}
}

func editTart(params map[string]string, w io.Writer) {
	if missingFields := checkHasFields([]string{"tart"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart edit-tart --tart <pushURL> [--name <name>] [--set-env \"<env-name>=<env-value>\"] [--delete-env <env-name>]")
		printMissingFields(missingFields, w)
		return
	}

	exists, tart := findTart(params["tart"])
	if !exists {
		fmt.Fprintln(w, "Err: A tart by that pushURL does not exist")
		return
	}

	if params["name"] != "" {
		tart.Name = params["name"]
	}

	if params["set-env"] != "" {
		tart.Env = setEnv(tart.Env, params["set-env"], "")
	}

	if params["delete-env"] != "" {
		tart.Env = setEnv(tart.Env, "", params["delete-env"])
	}

	tartmanager.Save(tart.PushURL, tart)
}

func setEnv(envList []string, envString, delString string) []string {
	key := strings.Split(envString, "=")[0]
	var output []string

	for _, envEntry := range envList {
		if (strings.Split(envEntry, "=")[0] == key) || strings.Split(envEntry, "=")[0] == delString {
			//no op
		} else {
			output = append(output, envEntry)
		}
	}
	if envString != "" {
		output = append(output, envString)
	}
	return output
}
