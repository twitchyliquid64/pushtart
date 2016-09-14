package cmd_registry

import "io"

var commands map[string]func(map[string]string, io.Writer)

//Register adds the given command and its function pointer to the registry.
func Register(cmd string, function func(map[string]string, io.Writer)) {
	if commands == nil {
		commands = map[string]func(map[string]string, io.Writer){}
	}
	commands[cmd] = function
}

//Command returns the function pointer for the specified command. false, nil is returned if the command does not exist.
func Command(cmd string) (ok bool, function func(map[string]string, io.Writer)) {
	function, ok = commands[cmd]
	return
}

//List returns a []string of all registered commands.
func List() []string {
	var output []string
	for key := range commands {
		output = append(output, key)
	}
	return output
}

//Init registers ssh-only commands with the registry.
func Init() {
	if commands == nil {
		commands = map[string]func(map[string]string, io.Writer){}
	}
	Register("exit", exit)
}

func exit(params map[string]string, w io.Writer) {
	//Caught / implemented elsewhere in the call chain - unreachable
}
