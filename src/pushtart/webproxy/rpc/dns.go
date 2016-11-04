package rpc

import (
	"pushtart/config"
	"pushtart/dnsserv"
	"pushtart/logging"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// DNSExtension represents the authenticated RPC server for DNSSERV-related methods, available via a webproxy URI.
type DNSExtension int

// ListRecordsResult is the return value of List().
type ListRecordsResult struct {
	A map[string]DNSRecord
}

// DNSRecord represents a DNS zone that pushtart serves.
type DNSRecord struct {
	Address string
	TTL     int
}

// List constructs a ListRecordsResult and returns it to the RPC caller. The structure contains information about
// the DNS records currently set in pushtart.
func (e *DNSExtension) List(arg map[string]string, result *ListRecordsResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg["APIKey"]); ok {
		logging.Info("rpc", "["+serviceName+"] List()")
	} else {
		logging.Warning("rpc", "Invalid auth for List()")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	result.A = make(map[string]DNSRecord)
	for name, obj := range config.All().DNS.ARecord {
		result.A[name] = DNSRecord{
			Address: obj.Address,
			TTL:     int(obj.TTL),
		}
	}
	return nil
}

// SetAArgs specifes the parameters passed to the DNSExtension.SetA RPC.
type SetAArgs struct {
	APIKey  string
	Domain  string
	Address string
	TTL     int
}

// SetA sets an A record for a specific domain, such that the DNS server will respond with the given address and record TTL.
func (e *DNSExtension) SetA(arg SetAArgs, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] SetA("+arg.Domain+" => "+arg.Address+")")
	} else {
		logging.Warning("rpc", "Invalid auth for SetA("+arg.Domain+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if config.All().DNS.ARecord == nil {
		config.All().DNS.ARecord = map[string]config.ARecord{}
	}

	config.All().DNS.ARecord[dnsserv.SanitizeDomain(arg.Domain)] = config.ARecord{
		Address: arg.Address,
		TTL:     uint32(arg.TTL),
	}
	config.Flush()
	result.Success = true
	return nil
}

// DeleteA deletes any DNS A records for the given domain.
func (e *DNSExtension) DeleteA(arg SetAArgs, result *ArbitrarySuccessResult) error {
	var serviceName string
	var ok bool
	if serviceName, ok = checkAuth(arg.APIKey); ok {
		logging.Info("rpc", "["+serviceName+"] DeleteA("+arg.Domain+")")
	} else {
		logging.Warning("rpc", "Invalid auth for DeleteA("+arg.Domain+")")
		return jsonrpc2.NewError(403, "Invalid API key")
	}

	if config.All().DNS.ARecord != nil {
		delete(config.All().DNS.ARecord, dnsserv.SanitizeDomain(arg.Domain))
	}
	config.Flush()
	result.Success = true
	return nil
}
