package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"strconv"
)

type HTTPConnectivityResult struct {
	resultCode  int
	redirectURL string
}

// CheckHTTPConnectivity Check basic HTTP connectivity to the domain
func CheckHTTPConnectivity(domain string, ip string) HTTPConnectivityResult {
	req := "GET http://" + ip + " HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n"
	if config.ProxyUsername != "" {
		req += "Proxy-Authorization: " + util.ParseAuth() + "\r\n"
	}
	req += "\r\n"

	resultCode, _, redirectURL, err := SendHTTPRequest(domain, ip, 80, req)
	if resultCode == -10 {
		util.PrintError(config.ProxyIP, domain, err)
	}

	util.PrintInfo(domain, domain+" returned "+strconv.Itoa(resultCode))

	return HTTPConnectivityResult{
		resultCode:  resultCode,
		redirectURL: redirectURL,
	}
}
