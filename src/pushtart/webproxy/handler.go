package webproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"pushtart/config"
	"pushtart/logging"
	"strconv"
	"strings"
	"time"
)

func main(w http.ResponseWriter, r *http.Request) {
	host := trimHostFieldToJustHostname(r.Host)
	if host == config.All().Web.DefaultDomain {
		internalsRouter.ServeHTTP(w, r)
	} else if isKnownVirtualDomain(host) {
		proxyRequestViaNetwork(config.All().Web.DomainProxies[host], w, r)
	} else {
		logging.Warning("httpproxy-main", "Request recieved for unknown virtual domain: "+host)
		internalsRouter.ServeHTTP(w, r)
	}
}

func isKnownVirtualDomain(host string) bool {
	_, ok := config.All().Web.DomainProxies[host]
	return ok
}

func trimHostFieldToJustHostname(hostField string) string {
	spl := strings.Split(hostField, ":")
	if len(spl) < 2 {
		return hostField
	}
	return spl[0]
}

func proxyRequestViaNetwork(proxyEntry config.DomainProxy, w http.ResponseWriter, r *http.Request) {
	if config.All().Web.LogAllProxies {
		logging.Info("httpproxy-main", "Proxying request "+r.Host+" -> "+proxyEntry.TargetHost+":"+strconv.Itoa(proxyEntry.TargetPort)+r.URL.Path)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = proxyEntry.TargetScheme
		req.URL.Path = r.URL.Path
		req.URL.Host = proxyEntry.TargetHost + ":" + strconv.Itoa(proxyEntry.TargetPort)
		req.Host = proxyEntry.TargetHost
	}

	prox := httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return http.ProxyFromEnvironment(req)
			},
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial(network, addr)
				if err != nil {
					logging.Warning("httpproxy-main", "Error connecting to backend: "+err.Error()+" ("+r.Host+")")
				}
				return conn, err
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	prox.ServeHTTP(w, r)
}
