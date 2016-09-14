package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

//Info writes a log line to the console + backlog of type info.
func Info(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "I", cyan(), content...))
	publishLogMessage("I", module, content...)
}

//Warning writes a log line to the console + backlog of type warning.
func Warning(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "W", yellow(), content...))
	publishLogMessage("W", module, content...)
}

//Error writes a log line to the console + backlog of type error.
func Error(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "E", red(), content...))
	publishLogMessage("E", module, content...)
}

//Fatal writes a log line to the console + backlog of type fatal, then terminates the program.
func Fatal(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "F", red(), content...))
	publishLogMessage("F", module, content...)
	os.Exit(1)
}

func formatLogPrefix(module, messagePrefix, prefixColor string, content ...interface{}) string {
	c := fmt.Sprint(content...)
	module = strings.ToUpper(module)
	if module != "" {
		return prefixColor + "[" + messagePrefix + "] " + blue() + "[" + module + "] " + clear() + c
	}
	return prefixColor + "[" + messagePrefix + "] " + clear() + c
}

var logSync sync.Mutex

func writeLogLine(inp string) {
	logSync.Lock()
	defer logSync.Unlock()
	log.Println(inp)
}

func publishLogMessage(msgType, module string, content ...interface{}) {
	c := fmt.Sprint(content...)
	publishMessage(module, msgType, c)
	addToBacklog(module, msgType, c)
}
