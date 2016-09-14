package main

import (
	"fmt"
	"io"
	"pushtart/constants"
	"pushtart/util"
	"strings"
)

func parseCommands(input []string) map[string]string {
	return defaultParams(util.ParseCommands(input))
}

func defaultParams(input map[string]string) map[string]string {
	if _, ok := input["config"]; !ok {
		input["config"] = constants.DefaultConfigFileName
	}
	return input
}

// checkHasFields returns a nil slice if all the fields in needFields are in input.
// otherwise, it returns the missing fields.
func checkHasFields(needFields []string, input map[string]string) []string {
	missingFields := []string{}

	for _, field := range needFields {
		if _, ok := input[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}

	return missingFields
}

func printMissingFields(missingFields []string, w io.Writer) {
	fmt.Fprintln(w, "Missing fields: "+strings.Join(missingFields, ","))
}
