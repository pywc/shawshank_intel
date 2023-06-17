package http_tester

type HTTPResult struct {
	connectivity         HTTPConnectivityResult
	headerHost           []FilteredHTTP
	htmlTitle            HTTPConnectivityResult
	htmlTokens           []FilteredHTTP
	hostnamePadding      []FilteredHTTP
	getCapitalize        []FilteredHTTP
	getRemove            []FilteredHTTP
	getAlternate         []FilteredHTTP
	httpCapitalize       []FilteredHTTP
	httpRemove           []FilteredHTTP
	httpAlternate        []FilteredHTTP
	hostCapitalize       []FilteredHTTP
	hostRemove           []FilteredHTTP
	hostAlternate        []FilteredHTTP
	httpDelimiterRemove  []FilteredHTTP
	pathAlternate        []FilteredHTTP
	headerAlternate      []FilteredHTTP
	hostnameAlternate    []FilteredHTTP
	hostnameTLDAlternate []FilteredHTTP
	subdomainAlternate   []FilteredHTTP
}

func TestHTTP(ip string, domain string) HTTPResult {
	result := HTTPResult{}
	result.connectivity = CheckHTTPConnectivity(domain, ip)
	_, result.headerHost = CheckHTTPHeaderHost(domain)
	result.htmlTitle = CheckHTMLTitle(domain)
	result.htmlTokens, _ = CheckHTMLTokens(domain)
	result.hostnamePadding = CheckHostnamePadding(domain, ip)
	result.getCapitalize = CheckGetWordCapitalize(domain, ip)
	result.getRemove = CheckGetWordRemove(domain, ip)
	result.getAlternate = CheckGetWordAlternate(domain, ip)
	result.httpCapitalize = CheckHTTPWordCapitalize(domain, ip)
	result.httpRemove = CheckHTTPWordRemove(domain, ip)
	result.httpAlternate = CheckHTTPWordAlternate(domain, ip)
	result.hostCapitalize = CheckHostWordCapitalize(domain, ip)
	result.hostRemove = CheckHostWordRemove(domain, ip)
	result.hostAlternate = CheckHostWordAlternate(domain, ip)
	result.httpDelimiterRemove = CheckHTTPDelimiterWordRemove(domain, ip)
	result.pathAlternate = CheckPathAlternate(domain, ip)
	result.hostnameAlternate = CheckHostnameAlternate(domain, ip)
	result.hostnameTLDAlternate = CheckHostnameTLDAlternate(domain, ip)
	result.subdomainAlternate = CheckHostnameSubdomainAlternate(domain, ip)

	return result
}
