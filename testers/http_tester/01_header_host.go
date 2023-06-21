package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"net/url"
	"strconv"
	"strings"
)

// CheckHTTPHeaderHost Checks whether wildcard-based filtering is used based on HTTP Host field
// Sets 399 as result code if it does not contain the correct magic word
/*
	Output Chart
	====================
	0: success
	1: exact-match
	2: wildcard-based
*/
func CheckHTTPHeaderHost(domain string) (int, []FilteredHTTP) {
	testList := config.DomainWildcards(domain)
	filteredList := make([]FilteredHTTP, len(testList))

	for _, testDomain := range testList {
		reqBody := "magicWord=" + url.QueryEscape(config.MagicWord) + "\r\n"
		req := "POST / HTTP/1.1\r\n" +
			"Host: " + testDomain + "\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: " + strconv.Itoa(len(reqBody)) + "\r\n"
		if config.ProxyType == "https" {
			req = "POST http://" + util.ParseEcho() + " HTTP/1.1\r\n" +
				"Host: " + testDomain + "\r\n" +
				"Content-Type: application/x-www-form-urlencoded\r\n" +
				"Content-Length: " + strconv.Itoa(len(reqBody)) + "\r\n"
			if config.ProxyUsername != "" {
				req += "Proxy-Authorization: Basic " + util.ParseAuth() + "\r\n"
			}
		}
		req += "\r\n"
		req += reqBody

		resultCode, resp, redirectURL, err := SendHTTPRequest(config.EchoServerAddr,
			config.EchoServerAddr, config.EchoServerPort, req, "")

		if resultCode == 0 && !strings.Contains(resp, config.MagicWord) {
			resultCode = 399
			redirectURL = "unknown"
		} else if resultCode == 0 {
			util.PrintInfo(domain, "header host result for "+testDomain+": "+strconv.Itoa(resultCode))
			continue
		} else if resultCode == -10 {
			util.PrintError(domain, err)
			continue
		}

		filtered := FilteredHTTP{
			Component:   testDomain,
			ResultCode:  resultCode,
			RedirectURL: redirectURL,
		}

		util.PrintInfo(domain, "header host result for "+testDomain+": "+strconv.Itoa(resultCode))
		filteredList = append(filteredList, filtered)
	}

	if len(filteredList) > 1 {
		return 2, filteredList
	} else if len(filteredList) > 0 {
		return 1, filteredList
	}

	return 0, filteredList
}
