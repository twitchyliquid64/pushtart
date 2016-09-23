package dnsserv

import (
	"errors"
	"net"
	"pushtart/config"
	"pushtart/logging"

	"github.com/miekg/dns"
)

//ErrRecordDoesNotExist is returned if a request is made without forwarding enabled and we have no record of it.
var ErrRecordDoesNotExist = errors.New("No authority for record")

func getRecordIfSpecifiedA(domain string) (dns.RR, error) {
	sDomain := sanitizeDomain(domain)
	if config.All().DNS.ARecord == nil {
		return nil, ErrRecordDoesNotExist
	}
	if out, ok := config.All().DNS.ARecord[sDomain]; ok {
		return makeAAnswer(domain, out.Address, out.TTL), nil
	}
	return nil, ErrRecordDoesNotExist
}

func queryA(domain string) (dns.RR, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (A): "+err.Error())
		return nil, err
	}
	for _, ip := range ips {
		if ip.To4() == nil {
			continue
		}
		return makeAAnswer(domain, ip.To4().String(), 3600), nil
	}
	return nil, errors.New("No A records for given query")
}

func getRecord(q dns.Question) (rr dns.RR, err error) {
	logging.Info("dnsserv-main", "Recieved query for "+q.Name)
	switch q.Qtype {
	case dns.TypeA:
		rr, err = getRecordIfSpecifiedA(q.Name)
		if err != nil && config.All().DNS.AllowForwarding {
			rr, err = queryA(q.Name)
		}
	}
	return rr, err
}

func parseQuery(m *dns.Msg) {
	var rr dns.RR

	for _, q := range m.Question {
		if resultRR, e := getRecord(q); e == nil {
			if resultRR == nil {
				continue
			}
			rr = resultRR.(dns.RR)
			logging.Info("dnsserv-debug", rr.String())
			if rr.Header().Name == q.Name {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	if config.All().DNS.AllowForwarding {
		m.RecursionAvailable = true
	}
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	case dns.OpcodeUpdate:
		logging.Warning("dnsserv-main", "Ignoring DNS request with opcodeUpdate")
	}

	// See http://mkaczanowski.com/golang-build-dynamic-dns-service-go/ for Tsig crap

	w.WriteMsg(m)
}
