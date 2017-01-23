package webproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"pushtart/config"
	"pushtart/logging"
	"pushtart/user"
	"strconv"
	"strings"
	"time"
)

//routes all requests
func main(w http.ResponseWriter, r *http.Request) {

	host := trimHostFieldToJustHostname(r.Host)
	if host == config.All().Web.DefaultDomain {
		internalsRouter.ServeHTTP(w, r)
	} else if isKnownVirtualDomain(host) {
		if config.All().TLS.Enabled && r.TLS == nil && config.All().TLS.ForceRedirect {
			redir(w, r)
			return
		}
		proxyRequestViaNetwork(config.All().Web.DomainProxies[host], w, r)
	} else {
		logging.Warning("httpproxy-main", "Request received for unknown virtual domain: "+host)
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

func redir(w http.ResponseWriter, req *http.Request) {
	portStr := ""
	if strings.Contains(config.All().TLS.Listener, ":") { //if a port is specified in the listening address
		portStr = config.All().TLS.Listener[strings.Index(config.All().TLS.Listener, ":"):] //get it
	}
	if portStr == ":443" {
		portStr = "" //no point, port 443 is implicit anyway for HTTPS
	}

	newURL := "https://" + req.Host + portStr + req.RequestURI
	logging.Info("httpproxy-main", "Redirecting ", req.URL, " to ", newURL)
	http.Redirect(w, req, newURL, http.StatusMovedPermanently)
}

func proxyRequestViaNetwork(proxyEntry config.DomainProxy, w http.ResponseWriter, r *http.Request) {

	//see if request authorized
	if !authorized(proxyEntry, w, r) {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Pushtart:"+config.All().Name+"\"")
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	//fetch the user if they exist
	username, _, ok := r.BasicAuth()
	var usr config.User
	if ok && user.Exists(username) {
		usr = user.Get(username)
	}

	if config.All().Web.LogAllProxies {
		logging.Info("httpproxy-main", "Proxying request "+r.Host+" -> "+proxyEntry.TargetHost+":"+strconv.Itoa(proxyEntry.TargetPort)+r.URL.Path)
	}
	director := func(req *http.Request) {
		req.URL.Scheme = proxyEntry.TargetScheme
		req.URL.Path = r.URL.Path
		req.URL.Host = proxyEntry.TargetHost + ":" + strconv.Itoa(proxyEntry.TargetPort)
		req.Host = proxyEntry.TargetHost
		if usr.Password != "" {
			req.Header.Add("X-auth-username", username)
			req.Header.Add("X-auth-name", usr.Name)
		}
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

//Grants authorization if there are not auth rules, else if they have a valid basic auth that positively matches
//one of the rules.
func authorized(proxyEntry config.DomainProxy, w http.ResponseWriter, r *http.Request) bool {
	if len(proxyEntry.AuthRules) == 0 {
		return true
	}
	usr, pwd, ok := r.BasicAuth()

	//First, check if any DENY rules match
	for _, rule := range proxyEntry.AuthRules {
		if rule.RuleType == "USR_DENY" {
			if usr == rule.Username {
				return false
			}
		}
	}

	if !ok { //No basic Auth specified
		return false
	}
	if !user.CheckUserPasswordWeb(usr, pwd) { //Incorrect password for given auth
		return false
	}
	for _, rule := range proxyEntry.AuthRules {
		switch rule.RuleType {
		case "ALLOW_ANY_USER":
			return true
		case "USR_ALLOW":
			if rule.Username == usr {
				return true
			}
		}
	}
	return false
}
