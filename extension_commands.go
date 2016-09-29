package main

import (
	"fmt"
	"io"
	"pushtart/config"
	"pushtart/dnsserv"
	"strconv"
	"strings"
)

func extensionCommand(params map[string]string, w io.Writer, user string) {
	if params["operation"] == "show-config" {
		listDnsservOptions(w)
		listHttpproxyOptions(w)
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
	if strings.ToUpper(params["extension"]) == "HTTPPROXY" {
		httpproxyCommand(params, w)
	}
}

func listHttpproxyOptions(w io.Writer) {
	fmt.Fprintln(w, "HTTPProxy:")
	fmt.Fprintln(w, "\t Enabled = "+strconv.FormatBool(config.All().Web.Enabled))
	fmt.Fprintln(w, "\t Listener = '"+config.All().Web.Listener+"'")
	fmt.Fprintln(w, "\t DefaultDomain = '"+config.All().Web.DefaultDomain+"'")
}

func httpproxyCommand(params map[string]string, w io.Writer) {
	if params["operation"] == "enable" {
		config.All().Web.Enabled = true
	}
	if params["operation"] == "disable" {
		config.All().Web.Enabled = false
	}
	if params["operation"] == "set-listener" && params["listener"] != "" {
		config.All().Web.Listener = params["listener"]
	} else if params["operation"] == "set-listener" {
		fmt.Fprintln(w, "Err: Missing fields: listener")
		return
	}
	if params["operation"] == "set-default-domain" && params["domain"] != "" {
		config.All().Web.DefaultDomain = params["domain"]
	} else if params["operation"] == "set-default-domain" {
		fmt.Fprintln(w, "Err: Missing fields: domain")
		return
	}

	if params["operation"] == "set-domain-proxy" {
		if missingFields := checkHasFields([]string{"extension", "operation", "domain", "targetport", "targethost"}, params); len(missingFields) > 0 {
			fmt.Fprintln(w, "USAGE: pushtart extension --extension HTTPProxy --operation set-domain-proxy --domain <domain> --targethost <host> --targetport <port>")
			printMissingFields(missingFields, w)
			return
		}

		if config.All().Web.DomainProxies == nil {
			config.All().Web.DomainProxies = map[string]config.DomainProxy{}
		}
		port, err := strconv.Atoi(params["targetport"])
		if err != nil {
			fmt.Fprintln(w, "Err parsing port: "+err.Error())
			return
		}

		scheme := "http"
		if params["scheme"] != "" {
			scheme = params["scheme"]
		}

		config.All().Web.DomainProxies[strings.ToLower(params["domain"])] = config.DomainProxy{
			TargetHost:   params["targethost"],
			TargetPort:   port,
			TargetScheme: scheme,
		}
	}
	if params["operation"] == "delete-domain-proxy" {
		if params["domain"] == "" {
			fmt.Fprintln(w, "Err: domain not specified")
			return
		}
		if config.All().Web.DomainProxies != nil {
			delete(config.All().Web.DomainProxies, strings.ToLower(params["domain"]))
		}
	}

	config.Flush()
}

func listDnsservOptions(w io.Writer) {
	fmt.Fprintln(w, "DNSServ:")
	fmt.Fprintln(w, "\t Enabled = "+strconv.FormatBool(config.All().DNS.Enabled))
	fmt.Fprintln(w, "\t Listener = '"+config.All().DNS.Listener+"'")
	fmt.Fprintln(w, "\t Allow-recursion = "+strconv.FormatBool(config.All().DNS.AllowForwarding))
	fmt.Fprintln(w, "\t Cache size = "+strconv.Itoa(config.All().DNS.LookupCacheSize))
}

func dnsservCommand(params map[string]string, w io.Writer) {
	if params["operation"] == "set-record" {
		if missingFields := checkHasFields([]string{"extension", "operation", "type", "domain", "address", "ttl"}, params); len(missingFields) > 0 {
			fmt.Fprintln(w, "USAGE: pushtart extension --extension DNSServ --operation set-record --type <DNS-record-type> --domain <domain> --address <ip-address> --ttl <expiry-seconds>")
			printMissingFields(missingFields, w)
			return
		}
	}

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
	if params["operation"] == "set-record" {
		if strings.ToUpper(params["type"]) == "A" {
			if config.All().DNS.ARecord == nil {
				config.All().DNS.ARecord = map[string]config.ARecord{}
			}
			ttl, err := strconv.Atoi(params["ttl"])
			if err != nil {
				fmt.Fprintln(w, "Err parsing ttl: "+err.Error())
				return
			}

			config.All().DNS.ARecord[dnsserv.SanitizeDomain(params["domain"])] = config.ARecord{
				Address: params["address"],
				TTL:     uint32(ttl),
			}
		}
	}
	if params["operation"] == "delete-record" {
		if params["domain"] == "" {
			fmt.Fprintln(w, "Err: domain not specified")
			return
		}
		if config.All().DNS.ARecord != nil {
			delete(config.All().DNS.ARecord, dnsserv.SanitizeDomain(params["domain"]))
		}
	}

	config.Flush()
}
