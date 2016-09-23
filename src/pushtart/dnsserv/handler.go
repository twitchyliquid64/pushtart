package dnsserv

import (
	"errors"
	"pushtart/config"
	"pushtart/logging"

	"github.com/miekg/dns"
)

//ErrRecordDoesNotExist is returned if a request is made without forwarding enabled and we have no record of it.
var ErrRecordDoesNotExist = errors.New("No authority for record")

func getRecordIfSpecifiedA(domain string) ([]dns.RR, error) {
	sDomain := sanitizeDomain(domain)
	if config.All().DNS.ARecord == nil {
		return nil, ErrRecordDoesNotExist
	}
	if out, ok := config.All().DNS.ARecord[sDomain]; ok {
		return []dns.RR{makeAAnswer(domain, out.Address, out.TTL)}, nil
	}
	return nil, ErrRecordDoesNotExist
}
func getRecordIfSpecifiedAAAA(domain string) ([]dns.RR, error) {
	sDomain := sanitizeDomain(domain)
	if config.All().DNS.AAAARecord == nil {
		return nil, ErrRecordDoesNotExist
	}
	if out, ok := config.All().DNS.AAAARecord[sDomain]; ok {
		return []dns.RR{makeAAAAAnswer(domain, out.Address, out.TTL)}, nil
	}
	return nil, ErrRecordDoesNotExist
}

func getRecord(q dns.Question) (rr []dns.RR, err error) {
	logging.Info("dnsserv-main", "Recieved query for "+q.Name)
	switch q.Qtype {
	case dns.TypeA:
		rr, err = getRecordIfSpecifiedA(q.Name)
		if err != nil && config.All().DNS.AllowForwarding {
			rr, err = queryA(q.Name)
		}
	case dns.TypeAAAA:
		rr, err = getRecordIfSpecifiedAAAA(q.Name)
		if err != nil && config.All().DNS.AllowForwarding {
			rr, err = queryAAAA(q.Name)
		}
	case dns.TypeNS:
		if config.All().DNS.AllowForwarding {
			rr, err = queryNS(q.Name)
		}
	case dns.TypeMX:
		if config.All().DNS.AllowForwarding {
			rr, err = queryMX(q.Name)
		}
	case dns.TypeTXT:
		if config.All().DNS.AllowForwarding {
			rr, err = queryTXT(q.Name)
		}
	case dns.TypeCNAME:
		if config.All().DNS.AllowForwarding {
			rr, err = queryCNAME(q.Name)
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
			for _, rr = range resultRR {
				logging.Info("dnsserv-debug", rr.String())
				if rr.Header().Name == q.Name {
					m.Answer = append(m.Answer, rr)
				}
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
