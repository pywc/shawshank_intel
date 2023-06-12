package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"log"
)

// CheckHTTPConnectivity Check basic HTTP connectivity to the domain
func CheckHTTPConnectivity(domain string, ip string) (int, string) {
	req := "GET / HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n\r\n"
	
	resultCode, _, redirectURL, err := SendHTTPRequest(domain, ip, 80, req)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	return resultCode, redirectURL
}
