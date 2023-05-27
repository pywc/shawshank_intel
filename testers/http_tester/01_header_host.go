package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"net/url"
	"strings"
)

type FilteredHeaderHost struct {
	host        string
	resultCode  int
	redirectURL string
}

// CheckHTTPHeaderHost Checks whether wildcard-based filtering is used based on HTTP Host field
// Sets 399 as result code if it does not contain the correct magic word
func CheckHTTPHeaderHost(domain string) []FilteredHeaderHost {
	testList := config.HTTPHostWildcards(domain)
	filteredList := make([]FilteredHeaderHost, 0)

	for _, testDomain := range testList {
		req := "POST / HTTP/1.1\r\nHost: " + testDomain + "\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
		req += "magicWord=" + url.QueryEscape(config.MagicWord)
		resultCode, resp, redirectURL := SendHTTPRequest(config.EchoServerAddr, config.EchoServerAddr, config.EchoServerPort, req)

		if resultCode == 0 && !strings.Contains(resp, config.MagicWord) {
			resultCode = 399
			redirectURL = "unknown"
		} else if resultCode == 0 {
			continue
		}

		filtered := FilteredHeaderHost{
			host:        testDomain,
			resultCode:  resultCode,
			redirectURL: redirectURL,
		}

		filteredList = append(filteredList, filtered)
	}

	return filteredList
}
