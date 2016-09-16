package util

import (
	"pushtart/logging"
	"strconv"
	"strings"
	"unicode"
)

//ParseCommands parses the given command strings - for each --<key> value it adds it to the returned map as map[key] = value.
func ParseCommands(input []string) map[string]string {
	out := map[string]string{}
	for i := 0; i < len(input); i++ {
		if strings.HasPrefix(input[i], "--") && len(input[i]) > 2 && (i+1) < len(input) {
			out[input[i][2:]] = input[i+1]
			i++
		}
	}
	return out
}

//BestPrefixMatch returns the given option of which inputStr is a subset (matching from the beginning of the string).
//If there is more than one such match or there are no matches, an empty string is returned.
func BestPrefixMatch(inputStr string, options []string) string {
	//iterate all the options, if there is one suffix match then use it.
	matches := 0
	lastMatch := ""
	for _, option := range options {
		if strings.HasPrefix(option, inputStr) {
			matches++
			lastMatch = option
		}
	}

	if matches == 1 {
		return lastMatch
	}
	return ""
}

//TokeniseCommandString takes the given string and parses it into separate string elements, following the same semantics for tokenising command parameters on the unix
//command line.
func TokeniseCommandString(input string) []string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	p := strings.FieldsFunc(input, f)
	for i := 0; i < len(p); i++ {
		if strings.HasPrefix(p[i], "\"") && strings.HasSuffix(p[i], "\"") {
			t, err := strconv.Unquote(p[i])
			if err == nil {
				p[i] = t
			} else {
				logging.Error("util-cmd-tokenise", err.Error())
			}
		}
	}

	return p
}
