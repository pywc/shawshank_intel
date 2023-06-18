package util

import (
	"github.com/pywc/shawshank_intel/testers/dns_tester"
	"github.com/pywc/shawshank_intel/testers/http_tester"
	"github.com/pywc/shawshank_intel/testers/https_tester"
	"github.com/pywc/shawshank_intel/testers/quic_tester"
)

type TestReport struct {
	country    string
	proxyIP    string
	hostDomain string
	hostIP     string
	dns        dns_tester.DNSResult
	http       http_tester.HTTPResult
	https      https_tester.HTTPSResult
	quic       quic_tester.QUICResult
}

func TestDomain(country string, proxyIP string, domain string) {
	hostIP, err := resolveIPLocally(domain)
	if err != nil {
		PrintError(proxyIP, domain, err)
		return
	}

	report := TestReport{
		country:    country,
		proxyIP:    proxyIP,
		hostDomain: domain,
		hostIP:     hostIP,
		dns:        dns_tester.TestDNS("", domain),
		http:       http_tester.TestHTTP("", domain),
		https:      https_tester.TestHTTPS("", domain),
		quic:       quic_tester.TestQUIC("", domain),
	}

	err = saveToDB(report)
	if err != nil {
		PrintError(proxyIP, domain, err)
	}
}

func saveToDB(report TestReport) error {
	return nil
}
