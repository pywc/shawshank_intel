package https_tester

import (
	utls "github.com/refraction-networking/utls"
	"log"
)

// CheckHTTPSConnectivity Check basic HTTPS connectivity to the domain
func CheckHTTPSConnectivity(domain string, ip string) int {
	// request configuration
	req := "GET / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"
	utlsConfig := utls.Config{
		ServerName: domain,
	}

	resultCode, _, err := SendHTTPSRequest(domain, ip, 443, req, &utlsConfig)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	return resultCode
}
