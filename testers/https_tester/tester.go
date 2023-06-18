package https_tester

import "github.com/pywc/shawshank_intel/util"

type HTTPSResult struct {
	Connectivity          int             `json:"connectivity"`
	SNI                   FilteredSNI     `json:"SNI"`
	ESNI                  int             `json:"ESNI,omitempty"`
	Certificate           int             `json:"certificate,omitempty"`
	BlindTLS              int             `json:"blindTLS,omitempty"`
	Splithello            int             `json:"splithello,omitempty"`
	SNIPadding            []FilteredHTTPS `json:"sni_padding,omitempty"`
	MinVerAlternate       []FilteredHTTPS `json:"min_ver_alternate,omitempty"`
	MaxVerAlternate       []FilteredHTTPS `json:"max_ver_alternate,omitempty"`
	CipherSuiteAlternate  []FilteredHTTPS `json:"cipher_suite_alternate,omitempty"`
	ClientCertAlternate   []FilteredHTTPS `json:"client_cert_alternate,omitempty"`
	SNITLDAlternate       []FilteredHTTPS `json:"sni_tld_alternate,omitempty"`
	SNISubdomainAlternate []FilteredHTTPS `json:"sni_subdomain_alternate,omitempty"`
}

func TestHTTPS(ip string, domain string) HTTPSResult {
	util.PrintInfo(domain, "testing HTTPS...")
	result := HTTPSResult{}
	result.Connectivity = CheckHTTPSConnectivity(domain, ip)
	if result.Connectivity == 0 {
		return result
	}

	result.ESNI = CheckESNI()
	result.Certificate = -5
	result.BlindTLS, _ = CheckTLS12Resumption(domain, ip)
	result.SNIPadding = CheckServerNamePadding(domain, ip)
	result.MinVerAlternate = CheckMinVersionAlternate(domain, ip)
	result.MaxVerAlternate = CheckMaxVersionAlternate(domain, ip)
	result.CipherSuiteAlternate = CheckCipherSuiteAlternate(domain, ip)
	result.ClientCertAlternate = CheckClientCertAlternate(domain, ip)
	result.SNITLDAlternate = CheckServernameTLDAlternate(domain, ip)
	result.SNISubdomainAlternate = CheckServernameSubdomainAlternate(domain, ip)

	return result
}
