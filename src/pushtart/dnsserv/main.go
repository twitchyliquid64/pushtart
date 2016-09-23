package dnsserv

import (
	"pushtart/config"
	"pushtart/logging"

	"github.com/miekg/dns"
)

func start() {
	server := &dns.Server{Addr: config.All().DNS.Listener, Net: "udp"}

	// Attach request handler func
	dns.HandleFunc(".", handleDNSRequest)

	//server.TsigSecret = map[string]string{name: secret}

	logging.Info("dnsserv-init", "Started")
	err := server.ListenAndServe()
	if err != nil {
		logging.Error("dnsserv-main", "ListenAndServe() err: "+err.Error())
	}
	defer server.Shutdown()
}

// Init is called by the main function to start the server - server will not be started if the DNS subsystem is disabled in configuration.
func Init() {
	if config.All().DNS.Enabled {

		err := initCache()
		if err != nil {
			logging.Error("dnsserv-init", "Error initializing Lookup Cache: "+err.Error())
		}

		go start()
	} else {
		logging.Info("dnsserv-init", "DNS is disabled - skipping init")
	}
}
