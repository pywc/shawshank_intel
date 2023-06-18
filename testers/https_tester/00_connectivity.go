package https_tester

import (
	"fmt"
	"github.com/pywc/shawshank_intel/config"
	utls "github.com/refraction-networking/utls"
	"log"
)

// CheckHTTPSConnectivity Check basic HTTPS Connectivity to the domain
func CheckHTTPSConnectivity(domain string, ip string) int {
	// request configuration
	req := "GET / HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n\r\n"

	utlsConfig := utls.Config{
		ServerName: domain,
	}

	resultCode, resp, err := SendHTTPSRequest(domain, ip, 443, req, &utlsConfig)
	fmt.Println(resp)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	return resultCode
}
