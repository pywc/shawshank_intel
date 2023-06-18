package dns_tester

type DNSResult struct {
	PublicDNS    int `json:"public_dns,omitempty"`
	DnsOverTLS   int `json:"dns_over_tls,omitempty"`
	DnsWhitelist int `json:"dns_whitelist,omitempty"`
}

func TestDNS(ip string, domain string) DNSResult {
	return DNSResult{
		PublicDNS:    -5,
		DnsOverTLS:   -5,
		DnsWhitelist: -5,
	}
}
