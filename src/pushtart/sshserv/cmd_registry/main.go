package cmd_registry

import "io"

var commands map[string]func(map[string]string, io.Writer)

func Register(cmd string, function func(map[string]string, io.Writer)) {
	if commands == nil {
		commands = map[string]func(map[string]string, io.Writer){}
	}
	commands[cmd] = function
}

func Command(cmd string) (ok bool, function func(map[string]string, io.Writer)) {
	function, ok = commands[cmd]
	return
}
