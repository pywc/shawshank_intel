package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
)

type HTTPResult struct {
	Connectivity         HTTPConnectivityResult `json:"connectivity"`
	HeaderHost           []FilteredHTTP         `json:"header_host,omitempty"`
	HtmlTitle            []FilteredHTTP         `json:"html_title"`
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
	config.CurrentComponent = "http"
	util.PrintInfo(domain, "testing HTTP...")

	result := HTTPResult{}

	result.Connectivity = CheckHTTPConnectivity(domain, ip)
	if result.Connectivity.ResultCode <= 0 {
		return result
	}

	redirectHost := result.Connectivity.RedirectURL

	_, result.HeaderHost = CheckHTTPHeaderHost(domain)
	result.HtmlTitle = CheckHTMLTitle(domain)
	result.HtmlTokens, _ = CheckHTMLTokens(domain)

	util.PrintInfo(domain, "initiating cenfuzz drivers...")

	if config.ProxyType != "https" {
		result.HostnamePadding = CheckHostnamePadding(domain, ip, redirectHost)
		result.GetCapitalize = CheckGetWordCapitalize(domain, ip, redirectHost)
		result.GetRemove = CheckGetWordRemove(domain, ip, redirectHost)
		result.HttpCapitalize = CheckHTTPWordCapitalize(domain, ip, redirectHost)
		result.HttpRemove = CheckHTTPWordRemove(domain, ip, redirectHost)
		result.HttpAlternate = CheckHTTPWordAlternate(domain, ip, redirectHost)
		result.HostCapitalize = CheckHostWordCapitalize(domain, ip, redirectHost)
		result.HostRemove = CheckHostWordRemove(domain, ip, redirectHost)
		result.HostAlternate = CheckHostWordAlternate(domain, ip, redirectHost)
		result.HttpDelimiterRemove = CheckHTTPDelimiterWordRemove(domain, ip, redirectHost)
	}

	result.GetAlternate = CheckGetWordAlternate(domain, ip, redirectHost)
	result.HeaderAlternate = CheckHeaderAlternate(domain, ip, redirectHost)
	result.PathAlternate = CheckPathAlternate(domain, ip, redirectHost)

	// TODO: fix error
	// result.HostnameAlternate = CheckHostnameAlternate(domain, ip, redirectHost)
	result.HostnameTLDAlternate = CheckHostnameTLDAlternate(domain, ip, redirectHost)
	result.SubdomainAlternate = CheckHostnameSubdomainAlternate(domain, ip, redirectHost)

	return result
}
