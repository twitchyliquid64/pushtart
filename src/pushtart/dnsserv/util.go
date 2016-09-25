package dnsserv

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

//SanitizeDomain appends a trailing dot to the string if none already exists.
func SanitizeDomain(in string) string {
	if !strings.HasSuffix(in, ".") {
		return in + "."
	}
	return in
}

func makeMXAnswer(name, answerDomain string, Pref uint16, TTL uint32) dns.RR {
	name = SanitizeDomain(name)
	answerDomain = SanitizeDomain(answerDomain)

	r := new(dns.MX)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: TTL}
	r.Preference = Pref
	r.Mx = answerDomain
	return r
}

func makeAAnswer(name, addr string, TTL uint32) dns.RR {
	name = SanitizeDomain(name)

	r := new(dns.A)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: TTL}
	r.A = net.ParseIP(addr)
	return r
}

//TODO: Rename to someone less retarded
func makeAAAAAnswer(name, addr string, TTL uint32) dns.RR {
	name = SanitizeDomain(name)

	r := new(dns.AAAA)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: TTL}
	r.AAAA = net.ParseIP(addr)
	return r
}

func makeNSAnswer(name, host string, TTL uint32) dns.RR {
	name = SanitizeDomain(name)

	r := new(dns.NS)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: TTL}
	r.Ns = host
	return r
}

func makeTXTAnswer(name, txt string, TTL uint32) dns.RR {
	name = SanitizeDomain(name)

	r := new(dns.TXT)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: TTL}
	r.Txt = []string{txt}
	return r
}

func makeCNAMEAnswer(name, names string, TTL uint32) dns.RR {
	name = SanitizeDomain(name)

	r := new(dns.CNAME)
	r.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: TTL}
	r.Target = names
	return r
}
