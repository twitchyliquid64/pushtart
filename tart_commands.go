package main

import (
	"fmt"
	"io"
	"pushtart/config"
	"pushtart/tartmanager"
)

func listTarts(params map[string]string, w io.Writer) {
	for pushURL, tart := range config.All().Tarts {
		fmt.Fprint(w, tart.Name+" ("+pushURL+"): ")
		if tart.IsRunning {
			fmt.Fprintln(w, "Running (PID ", tart.PID, ")")
		} else {
			fmt.Fprintln(w, "Stopped.")
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
