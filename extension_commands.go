package main

import (
	"fmt"
	"io"
	"pushtart/config"
	"strconv"
	"strings"
)

func extensionCommand(params map[string]string, w io.Writer) {
	if params["operation"] == "show-config" {
		listDnsservOptions(w)
		return
	}

	if missingFields := checkHasFields([]string{"extension"}, params); len(missingFields) > 0 {
		fmt.Fprintln(w, "USAGE: pushtart extension --extension <extension>")
		printMissingFields(missingFields, w)
		return
	}

	if strings.ToUpper(params["extension"]) == "DNSSERV" {
		dnsservCommand(params, w)
	}
}

func listDnsservOptions(w io.Writer) {
	fmt.Fprintln(w, "DNSServ:")
	fmt.Fprintln(w, "\t Enabled = "+strconv.FormatBool(config.All().DNS.Enabled))
	fmt.Fprintln(w, "\t Listener = '"+config.All().DNS.Listener+"'")
	fmt.Fprintln(w, "\t Allow-recursion = "+strconv.FormatBool(config.All().DNS.AllowForwarding))
	fmt.Fprintln(w, "\t Cache size = "+strconv.Itoa(config.All().DNS.LookupCacheSize))
}

func dnsservCommand(params map[string]string, w io.Writer) {
	if params["cache-size"] != "" {
		cs, err := strconv.Atoi(params["cache-size"])
		if err != nil {
			fmt.Fprintln(w, "Err parsing cache-size: "+err.Error())
			return
		}
		config.All().DNS.LookupCacheSize = cs
	}

	if params["operation"] == "enable" {
		config.All().DNS.Enabled = true
	}
	if params["operation"] == "disable" {
		config.All().DNS.Enabled = false
	}
	if params["operation"] == "enable-recursion" {
		config.All().DNS.AllowForwarding = true
	}
	if params["operation"] == "disable-recursion" {
		config.All().DNS.AllowForwarding = false
	}

	config.Flush()
}
