package dnsserv

import (
	"errors"
	"net"
	"pushtart/config"
	"pushtart/logging"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

var lookupCache *lru.ARCCache
var cacheExpiryTime = 60 * 20

func initCache() error {
	var err error
	lookupCache, err = lru.NewARC(config.All().DNS.LookupCacheSize)
	return err
}

func queryA(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("A:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("A:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (A): "+err.Error())
		return nil, err
	}
	if len(ips) > 0 {
		var output []dns.RR
		for _, ip := range ips {
			if ip.To4() == nil {
				continue
			}
			output = append(output, makeAAnswer(domain, ip.To4().String(), 3600))
		}
		lookupCache.Add("A:"+domain, cacheCollection{Answers: output, Timestamp: time.Now()})
		return output, nil
	}
	return nil, errors.New("No A records for given query")
}

func queryAAAA(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("AAAA:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("AAAA:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (AAAA): "+err.Error())
		return nil, err
	}
	var output []dns.RR

	if len(ips) > 0 {
		for _, ip := range ips {
			if ip.To4() != nil {
				continue
			}
			output = append(output, makeAAAAAnswer(domain, ip.To16().String(), 3600))
		}
		lookupCache.Add("AAAA:"+domain, cacheCollection{Answers: output, Timestamp: time.Now()})
		return output, nil
	}

	return nil, errors.New("No AAAA records for given query")
}

func queryNS(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("NS:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("NS:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	nss, err := net.LookupNS(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (NS): "+err.Error())
		return nil, err
	}
	if len(nss) > 0 {
		var output []dns.RR
		for _, ns := range nss {
			output = append(output, makeNSAnswer(domain, ns.Host, 3600))
		}
		lookupCache.Add("NS:"+domain, cacheCollection{Answers: output, Timestamp: time.Now()})
		return output, nil
	}
	return nil, errors.New("No NS records for given query")
}

func queryMX(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("MX:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("MX:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	nss, err := net.LookupMX(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (MX): "+err.Error())
		return nil, err
	}
	if len(nss) > 0 {
		var output []dns.RR
		for _, mx := range nss {
			output = append(output, makeMXAnswer(domain, mx.Host, mx.Pref, 3600))
		}
		lookupCache.Add("MX:"+domain, cacheCollection{Answers: output, Timestamp: time.Now()})
		return output, nil
	}
	return nil, errors.New("No MX records for given query")
}

func queryTXT(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("TXT:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("TXT:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	txts, err := net.LookupTXT(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (TXT): "+err.Error())
		return nil, err
	}
	if len(txts) > 0 {
		var output []dns.RR
		for _, txt := range txts {
			output = append(output, makeTXTAnswer(domain, txt, 3600))
		}
		lookupCache.Add("TXT:"+domain, cacheCollection{Answers: output, Timestamp: time.Now()})
		return output, nil
	}
	return nil, errors.New("No TXT records for given query")
}

func queryCNAME(domain string) ([]dns.RR, error) {
	cacheVal, cacheHit := lookupCache.Get("CNAME:" + domain)
	if cacheHit {
		if cacheVal.(cacheCollection).Timestamp.Before(time.Now().Add(-1 * time.Second * time.Duration(cacheExpiryTime))) {
			lookupCache.Remove("CNAME:" + domain)
		} else {
			logging.Info("dnsserv-main", "Cache hit")
			return cacheVal.(cacheCollection).Answers, nil
		}
	}

	names, err := net.LookupCNAME(domain)
	if err != nil {
		logging.Warning("dnsserv-main", "Lookup failure (CNAME): "+err.Error())
		return nil, err
	}

	ans := makeCNAMEAnswer(domain, names, 3600)
	lookupCache.Add("CNAME:"+domain, cacheValue{Answer: ans, Timestamp: time.Now()})
	return []dns.RR{ans}, nil
}

type cacheValue struct {
	Answer    dns.RR
	Timestamp time.Time
}
type cacheCollection struct {
	Answers   []dns.RR
	Timestamp time.Time
}
