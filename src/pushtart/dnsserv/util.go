package dnsserv

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func sanitizeDomain(in string) string {
	if !strings.HasSuffix(in, ".") {
		return in + "."
	}
	return in
}

func makeMXAnswer(name, answerDomain string, TTL uint32) dns.RR {
	name = sanitizeDomain(name)
	answerDomain = sanitizeDomain(answerDomain)

	r := new(dns.MX)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: TTL}
	r.Preference = 10
	r.Mx = answerDomain
	return r
}

func makeAAnswer(name, addr string, TTL uint32) dns.RR {
	name = sanitizeDomain(name)

	r := new(dns.A)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: TTL}
	r.A = net.ParseIP(addr)
	return r
}
