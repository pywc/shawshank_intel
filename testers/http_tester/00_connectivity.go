package http_tester

// CheckHTTPConnectivity Check basic HTTP connectivity to the domain
func CheckHTTPConnectivity(domain string, ip string) (int, string) {
	req := "POST / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"
	resultCode, _, redirectURL := SendHTTPRequest(domain, ip, 80, req)

	return resultCode, redirectURL
}
