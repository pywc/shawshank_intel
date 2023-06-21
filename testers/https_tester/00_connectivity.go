package https_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	"log"
	"strconv"
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

	resultCode, err := SendHTTPSRequestNormally(domain, ip, 443, req, &utlsConfig)
	if resultCode != 0 {
		return resultCode
	}

	resultCode, _, err = SendHTTPSRequest(domain, ip, 443, req, &utlsConfig)
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	util.PrintInfo(domain, domain+" returned "+strconv.Itoa(resultCode))

	return resultCode
}
