package webproxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"pushtart/config"
	"pushtart/logging"
)

var internalsRouter = http.NewServeMux()

// Init is called by the main function to start the server - server will not be started if the Web subsystem is disabled in configuration.
func Init() {
	if config.All().Web.Enabled {

		initRoutes()
		go start()
	} else {
		logging.Info("httpproxy-init", "HTTPProxy is disabled - skipping init")
	}
}

func initRoutes() {
	internalsRouter.HandleFunc("/health", health)
	internalsRouter.HandleFunc("/status", statusPage)
	internalsRouter.HandleFunc("/pubrpc", pubRPCInit().ServeHTTP)
	http.HandleFunc("/", main)
}

func health(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Ok")
}

func start() {
	logging.Info("httpproxy-init", "Initialising HTTP server on ", config.All().Web.Listener)

	if config.All().TLS.Enabled {
		go http.ListenAndServe(config.All().Web.Listener, nil)
		logging.Info("httpproxy-init", "Initialising HTTPS server on ", config.All().TLS.Listener)
		listener, err := tls.Listen("tcp", config.All().TLS.Listener, config.TLS())
		if err != nil {
			logging.Info("httpproxy-init", "Error creating listener: "+err.Error())
			return
		}

		http.Serve(listener, nil)
	} else {
		http.ListenAndServe(config.All().Web.Listener, nil)
	}
}
