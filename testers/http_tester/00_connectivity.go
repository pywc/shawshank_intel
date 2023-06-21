package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"strconv"
)

type HTTPConnectivityResult struct {
	ResultCode  int    `json:"result_code,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// CheckHTTPConnectivity Check basic HTTP connectivity to the domain
func CheckHTTPConnectivity(domain string, ip string) HTTPConnectivityResult {
	resultCode, redirectHost, err := SendHTTPRequestNormally(domain, ip, 80)
	if err != nil {
		util.PrintError(domain, err)

		return HTTPConnectivityResult{
			ResultCode:  resultCode,
			RedirectURL: redirectHost,
		}
	} else if resultCode >= 400 {
		return HTTPConnectivityResult{
			ResultCode:  -3,
			RedirectURL: redirectHost,
		}
	}

	req := "GET / HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n\r\n"
	if config.ProxyType == "https" {
		req = "GET http://" + ip + " HTTP/1.1\r\n" +
			"Host: " + domain + "\r\n" +
			"Accept: */*\r\n" +
			"User-Agent: " + config.UserAgent + "\r\n"
		if config.ProxyUsername != "" {
			req += "Proxy-Authorization: Basic " + util.ParseAuth() + "\r\n"
		}
	}

	req += "\r\n"

	resultCode, _, redirectURL, err := SendHTTPRequest(domain, ip, 80, req, redirectHost)
	if resultCode == -10 {
		util.PrintError(domain, err)
	}

	util.PrintInfo(domain, domain+" returned "+strconv.Itoa(resultCode)+" "+redirectURL)

	return HTTPConnectivityResult{
		ResultCode:  resultCode,
		RedirectURL: redirectURL,
	}
}
