package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"log"
	"net/url"
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
		req := "POST / HTTP/1.1\r\nHost: " + testDomain + "\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
		req += "magicWord=" + url.QueryEscape(config.MagicWord)
		resultCode, resp, redirectURL, err := SendHTTPRequest(config.EchoServerAddr, config.EchoServerAddr, config.EchoServerPort, req)
		if resultCode == 0 {
			continue
		} else if resultCode == -10 {
			log.Println("[*] Error - " + domain + " - " + err.Error())
		} else if resultCode == 0 && !strings.Contains(resp, config.MagicWord) {
			resultCode = 399
			redirectURL = "unknown"
		}

		filtered := FilteredHTTP{
			component:   testDomain,
			resultCode:  resultCode,
			redirectURL: redirectURL,
		}

		filteredList = append(filteredList, filtered)
	}

	if len(filteredList) > 1 {
		return 2, filteredList
	} else if len(filteredList) > 0 {
		return 1, filteredList
	}

	return 0, filteredList
}
