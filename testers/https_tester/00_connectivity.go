package https_tester

import (
	"crypto/tls"
	"log"
)

// CheckHTTPSConnectivity
/*
	-3: unhandled error
	-2: proxy error
	-1: not accessible in the US
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
*/
func CheckHTTPSConnectivity(domain string, ip string) int {
	// request configuration
	req := "GET / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"
	tlsConfig := tls.Config{
		ServerName: domain,
	}

	resultCode, _, err := SendHTTPSRequest(domain, ip, 443, req, &tlsConfig)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	return resultCode
}
