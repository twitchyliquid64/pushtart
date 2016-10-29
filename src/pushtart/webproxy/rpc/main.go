package rpc

import (
	"pushtart/config"
	"pushtart/logging"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

//AuthenticationArgument is the default argument for proceedures which require no arguments but authentication.
type AuthenticationArgument struct {
	APIKey string
}

//ListTartsResult represents the result of a successful ListTarts RPC.
type ListTartsResult struct {
	Tarts map[string]config.Tart
}

//Service represents the authenticated RPC server available via a webproxy URI.
type Service int

//ListTarts RPC returns a list of tarts.
func (t *Service) ListTarts(arg *AuthenticationArgument, result *ListTartsResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] ListTarts()")
	} else {
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	result.Tarts = map[string]config.Tart{}
	for name, t := range config.All().Tarts {
		result.Tarts[name] = t
	}
	return nil
}

func checkAuth(key string) (service string, ok bool) {
	for _, entry := range config.All().APIKeys {
		if entry.Key == key {
			return entry.Service, true
		}
	}
	return "", false
}
