package https_tester

type HTTPSResult struct {
	connectivity          int
	sni                   FilteredSNI
	esni                  int
	certificate           int
	blindtls              int
	splithello            int
	sniPadding            []FilteredHTTPS
	minVerAlternate       []FilteredHTTPS
	maxVerAlternate       []FilteredHTTPS
	cipherSuiteAlternate  []FilteredHTTPS
	clientCertAlternate   []FilteredHTTPS
	sniTLDAlternate       []FilteredHTTPS
	sniSubdomainAlternate []FilteredHTTPS
}

func TestHTTPS(ip string, domain string) HTTPSResult {
	result := HTTPSResult{}
	result.connectivity = CheckHTTPSConnectivity(domain, ip)
	result.esni = CheckESNI()
	result.certificate = -5
	result.blindtls, _ = CheckTLS12Resumption(domain, ip)
	result.sniPadding = CheckServerNamePadding(domain, ip)
	result.minVerAlternate = CheckMinVersionAlternate(domain, ip)
	result.maxVerAlternate = CheckMaxVersionAlternate(domain, ip)
	result.cipherSuiteAlternate = CheckCipherSuiteAlternate(domain, ip)
	result.clientCertAlternate = CheckClientCertAlternate(domain, ip)
	result.sniTLDAlternate = CheckServernameTLDAlternate(domain, ip)
	result.sniSubdomainAlternate = CheckServernameSubdomainAlternate(domain, ip)

	return result
}
