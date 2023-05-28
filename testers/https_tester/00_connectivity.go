package https_tester

import (
	"crypto/tls"
	"log"
)

// CheckHTTPSConnectivity Check basic HTTPS connectivity to the domain
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
