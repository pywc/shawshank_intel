package https_tester

import (
	"github.com/pywc/shawshank_intel/config"
	utls "github.com/refraction-networking/utls"
	"log"
)

type FilteredSNI struct {
	sni        string
	resultCode int
}

func CheckSNI(domain string, ip string) (int, []FilteredSNI) {
	testList := config.DomainWildcards(domain)
	filteredList := make([]FilteredSNI, len(testList))

	for _, testDomain := range testList {
		req := "GET / HTTP/1.1\r\n" +
			"Host: " + domain + "\r\n" +
			"Accept: */*\r\n" +
			"User-Agent: " + config.UserAgent + "\r\n\r\n"

		utlsConfig := utls.Config{
			InsecureSkipVerify: true,
			ServerName:         testDomain,
		}
		resultCode, _, err := SendHTTPSRequest(config.DummyServerDomain, config.DummyServerIP, 443, req, &utlsConfig)

		if resultCode == 0 {
			continue
		} else if resultCode == -10 {
			log.Println("[*] Error - " + domain + " - " + err.Error())
		}

		filtered := FilteredSNI{
			sni:        testDomain,
			resultCode: resultCode,
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
