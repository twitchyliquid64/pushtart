package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func Info(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "I", cyan(), content...))
	publishLogMessage("I", module, content...)
}

func Warning(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "W", yellow(), content...))
	publishLogMessage("W", module, content...)
}

func Error(module string, content ...interface{}) {
	writeLogLine(formatLogPrefix(module, "E", red(), content...))
	publishLogMessage("E", module, content...)
}

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
	} else {
		return prefixColor + "[" + messagePrefix + "] " + clear() + c
	}
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
