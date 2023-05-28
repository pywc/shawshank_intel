package http_tester

import "log"

// CheckHTTPConnectivity Check basic HTTP connectivity to the domain
func CheckHTTPConnectivity(domain string, ip string) (int, string) {
	req := "POST / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"
	resultCode, _, redirectURL, err := SendHTTPRequest(domain, ip, 80, req)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	return resultCode, redirectURL
}
