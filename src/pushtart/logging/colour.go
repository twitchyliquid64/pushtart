package logging

import (
	"runtime"
)

var enableColors = false

func init() {
	if runtime.GOOS == "windows" {
		enableColors = false
	} else {
		enableColors = true
	}
}

func green() string {
	return valIfColoursEnabled("\033[32;1m")
}
func yellow() string {
	return valIfColoursEnabled("\033[33;1m")
}
func blue() string {
	return valIfColoursEnabled("\033[34;1m")
}
func cyan() string {
	return valIfColoursEnabled("\033[36;1m")
}
func red() string {
	return valIfColoursEnabled("\033[31;1m")
}

func clear() string {
	return valIfColoursEnabled("\033[0m")
}

func valIfColoursEnabled(input string) string {
	if enableColors {
		return input
	}
	return ""
}
