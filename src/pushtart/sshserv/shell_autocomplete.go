package sshserv

import (
	"pushtart/sshserv/cmd_registry"
	"pushtart/tartmanager"
	"pushtart/user"
	"pushtart/util"
	"strings"
)

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

		case spl[len(spl)-2] == "--extension":
			if match := util.BestPrefixMatch(spl[len(spl)-1], []string{"DNSServ", "HTTPProxy"}); match != "" {
				newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
				return newLine, len(newLine), true
			}

		case spl[len(spl)-2] == "--operation":
			if spl[0] == "extension" {
				newLine, cursorPos, ok := tryMatchExtensionOperation(line, spl)
				if ok {
					return newLine, cursorPos, true
				}
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

func tryMatchExtensionOperation(line string, spl []string) (newLine string, cursorPos int, replace bool) {
	params := util.ParseCommands(util.TokeniseCommandString(line))
	if params["extension"] == "" {
		if match := util.BestPrefixMatch(spl[len(spl)-1], []string{"show-config"}); match != "" {
			newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
			return newLine, len(newLine), true
		}
	} else {
		operationPossibilities, ok := operationsByExtension[strings.ToUpper(params["extension"])]
		if ok {
			if match := util.BestPrefixMatch(spl[len(spl)-1], operationPossibilities); match != "" {
				newLine := strings.Join(spl[0:len(spl)-1], " ") + " " + match
				return newLine, len(newLine), true
			}
		}
	}
	return "", 0, false
}

var operationsByExtension = map[string][]string{
	"DNSSERV":   []string{"set-record", "delete-record", "enable", "enable-recursion", "disable", "disable-recursion"},
	"HTTPPROXY": []string{"enable", "disable", "set-listener", "set-default-domain", "set-domain-proxy", "delete-domain-proxy", "add-authorization-rule", "remove-authorization-rule"},
}

var commandParams = map[string][]string{
	"edit-user":         []string{"--username", "--password", "--name", "--allow-ssh-password"},
	"make-user":         []string{"--username", "--password", "--name", "--allow-ssh-password"},
	"start-tart":        []string{"--tart"},
	"stop-tart":         []string{"--tart"},
	"edit-tart":         []string{"--tart", "--name", "--set-env", "--delete-env", "--log-stdout"},
	"tart-restart-mode": []string{"--tart", "--enabled", "--lull-period"},
	"extension":         []string{"--extension", "--operation", "--domain", "--type"},
	"set-config-value":  []string{"--field", "--value"},
	"get-config-value":  []string{"--field"},
	"tart-add-owner":    []string{"--username", "--tart"},
	"tart-remove-owner": []string{"--username", "--tart"},
}
