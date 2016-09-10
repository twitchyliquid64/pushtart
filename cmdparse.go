package main

import (
  "pushtart/constants"
  "strings"
  "fmt"
)

func parseCommands(input []string)map[string]string{
  out := map[string]string{}
  for i := 0; i < len(input); i++{
    if strings.HasPrefix(input[i], "--") && len(input[i]) > 2 && (i+1) < len(input){
      out[input[i][2:]] = input[i+1]
      i++
    }
  }
  return defaultParams(out)
}


func defaultParams(input map[string]string)map[string]string{
  if _, ok := input["config"]; !ok {
    input["config"] = constants.DefaultConfigFileName
  }
  return input
}

// checkHasFields returns a nil slice if all the fields in needFields are in input.
// otherwise, it returns the missing fields.
func checkHasFields(needFields []string, input map[string]string)[]string{
  missingFields := []string{}

  for _, field := range needFields {
    if _, ok := input[field]; !ok {
      missingFields = append(missingFields, field)
    }
  }

  return missingFields
}


func printMissingFields(missingFields []string){
  fmt.Println("Missing fields: " + strings.Join(missingFields, ","))
}
