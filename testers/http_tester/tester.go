package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
)

type HTTPResult struct {
	Connectivity         HTTPConnectivityResult `json:"connectivity"`
	HeaderHost           []FilteredHTTP         `json:"header_host,omitempty"`
	HtmlTitle            HTTPConnectivityResult `json:"html_title"`
	HtmlTokens           []FilteredHTTP         `json:"html_tokens,omitempty"`
	HostnamePadding      []FilteredHTTP         `json:"hostname_padding,omitempty"`
	GetCapitalize        []FilteredHTTP         `json:"get_capitalize,omitempty"`
	GetRemove            []FilteredHTTP         `json:"get_remove,omitempty"`
	GetAlternate         []FilteredHTTP         `json:"get_alternate,omitempty"`
	HttpCapitalize       []FilteredHTTP         `json:"http_capitalize,omitempty"`
	HttpRemove           []FilteredHTTP         `json:"http_remove,omitempty"`
	HttpAlternate        []FilteredHTTP         `json:"http_alternate,omitempty"`
	HostCapitalize       []FilteredHTTP         `json:"host_capitalize,omitempty"`
	HostRemove           []FilteredHTTP         `json:"host_remove,omitempty"`
	HostAlternate        []FilteredHTTP         `json:"host_alternate,omitempty"`
	HttpDelimiterRemove  []FilteredHTTP         `json:"http_delimiter_remove,omitempty"`
	PathAlternate        []FilteredHTTP         `json:"path_alternate,omitempty"`
	HeaderAlternate      []FilteredHTTP         `json:"header_alternate,omitempty"`
	HostnameAlternate    []FilteredHTTP         `json:"hostname_alternate,omitempty"`
	HostnameTLDAlternate []FilteredHTTP         `json:"hostname_tld_alternate,omitempty"`
	SubdomainAlternate   []FilteredHTTP         `json:"subdomain_alternate,omitempty"`
}

func TestHTTP(ip string, domain string) HTTPResult {
	util.PrintInfo(domain, "testing HTTP...")
	result := HTTPResult{}
	result.Connectivity = CheckHTTPConnectivity(domain, ip)
	if result.Connectivity.resultCode == 0 || result.Connectivity.resultCode == -10 {
		return result
	}

	_, result.HeaderHost = CheckHTTPHeaderHost(domain)
	result.HtmlTitle = CheckHTMLTitle(domain)
	result.HtmlTokens, _ = CheckHTMLTokens(domain)

	util.PrintInfo(domain, "Initiating cenfuzz drivers...")

	if config.ProxyType != "https" {
		result.HostnamePadding = CheckHostnamePadding(domain, ip)
		result.GetCapitalize = CheckGetWordCapitalize(domain, ip)
		result.GetRemove = CheckGetWordRemove(domain, ip)
		result.GetAlternate = CheckGetWordAlternate(domain, ip)
		result.HttpCapitalize = CheckHTTPWordCapitalize(domain, ip)
		result.HttpRemove = CheckHTTPWordRemove(domain, ip)
		result.HttpAlternate = CheckHTTPWordAlternate(domain, ip)
		result.HostCapitalize = CheckHostWordCapitalize(domain, ip)
		result.HostRemove = CheckHostWordRemove(domain, ip)
		result.HostAlternate = CheckHostWordAlternate(domain, ip)
		result.HttpDelimiterRemove = CheckHTTPDelimiterWordRemove(domain, ip)
	}
	
	result.HeaderAlternate = CheckHeaderAlternate(domain, ip)
	result.PathAlternate = CheckPathAlternate(domain, ip)
	result.HostnameAlternate = CheckHostnameAlternate(domain, ip)
	result.HostnameTLDAlternate = CheckHostnameTLDAlternate(domain, ip)
	result.SubdomainAlternate = CheckHostnameSubdomainAlternate(domain, ip)

	return result
}
