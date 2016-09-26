package main

import (
	"errors"
	"fmt"
	"io"
	"pushtart/config"
	"pushtart/logging"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func logMsgs(params map[string]string, w io.Writer) {
	logMsgs := logging.GetBacklog()

	for _, msg := range logMsgs {
		fmt.Fprintln(w, time.Unix(msg.Created, 0).Format(time.ANSIC), "["+msg.Component+"]", msg.Message)
	}
}

func setConfigValue(params map[string]string, w io.Writer) {
	if missingFields := checkHasFields([]string{"field", "value"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart set-config-value --field <config-field> --value <new-value>")
		printMissingFields(missingFields, w)
		return
	}

	err := setVal(params["field"], params["value"], reflect.ValueOf(config.All()).Elem())
	if err != nil {
		fmt.Fprintln(w, "Err: "+err.Error())
	}
	config.Flush()
}

func getConfigValue(params map[string]string, w io.Writer) {
	if missingFields := checkHasFields([]string{"field"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart get-config-value --field <config-field>")
		printMissingFields(missingFields, w)
		return
	}

	out, err := getVal(params["field"], reflect.ValueOf(config.All()).Elem())
	if err != nil {
		fmt.Fprintln(w, "Err: "+err.Error())
		return
	}
	fmt.Fprintln(w, out)
}

func getVal(query string, val reflect.Value) (out string, err error) {
	spl := strings.Split(query, ".")
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		tag := typeField.Tag
		if len(spl) == 1 { //no other sections like DNS.Listener
			if typeField.Name == spl[0] && tag.Get("getConfigValue") != "block" {
				return valToStr(valueField), nil
			}
		} else {
			if typeField.Name == spl[0] {
				return getVal(strings.Join(spl[1:], "."), valueField)
			}
		}
	}

	return "", errors.New("Could not find field")
}

func setVal(query, newVal string, val reflect.Value) (err error) {
	spl := strings.Split(query, ".")
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		tag := typeField.Tag
		if len(spl) == 1 { //no other sections like DNS.Listener
			if typeField.Name == spl[0] && tag.Get("getConfigValue") != "block" {
				v, err := strToVal(newVal, valueField)
				if err != nil {
					return err
				}
				if !valueField.CanSet() {
					return errors.New("Cannot set field: " + typeField.Name)
				}
				valueField.Set(v)
				return nil
			}
		} else {
			if typeField.Name == spl[0] {
				return setVal(strings.Join(spl[1:], "."), newVal, valueField)
			}
		}
	}

	return errors.New("Could not find field")
}

func valToStr(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.String:
		return val.String()
	case reflect.Int64:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int:
		return strconv.Itoa(int(val.Int()))
	}
	return "?"
}

func strToVal(in string, template reflect.Value) (reflect.Value, error) {
	switch template.Kind() {
	case reflect.Bool:
		vb, err := strconv.ParseBool(in)
		return reflect.ValueOf(vb), err
	case reflect.String:
		return reflect.ValueOf(in), nil
	case reflect.Int64:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int:
		vi, err := strconv.Atoi(in)
		return reflect.ValueOf(vi), err
	}
	return reflect.ValueOf(nil), errors.New("Don't know how to process " + template.Kind().String())
}
