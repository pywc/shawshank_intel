package dns_tester

type DNSResult struct {
	publicDNS    int
	dnsOverTLS   int
	dnsWhitelist int
}

func TestDNS(ip string, domain string) DNSResult {
	return DNSResult{
		publicDNS:    -5,
		dnsOverTLS:   -5,
		dnsWhitelist: -5,
	}
}
