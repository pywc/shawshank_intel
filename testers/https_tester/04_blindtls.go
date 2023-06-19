package https_tester

import (
	"fmt"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	utls "github.com/refraction-networking/utls"
	url2 "net/url"
	"strings"
)

// CheckTLS12Resumption
func CheckTLS12Resumption(domain string, ip string) (int, error) {
	req := "GET / HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n\r\n"

	utlsConfig := utls.Config{
		ServerName: domain,
		MinVersion: utls.VersionTLS12,
		MaxVersion: utls.VersionTLS12,
	}

	conn, err := util.ConnectNormally(ip, 443)
	sess, err := util.GetNewTLSSession(conn, req, &utlsConfig)
	conn.Close()

	if err != nil {
		// unknown error
		return -10, err
	}

	conn, err = util.ConnectViaProxy(ip, 443, "https")
	resp, err := util.ResumeTLSSession(conn, req, *sess)
	conn.Close()

	// check tcp errors
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			// connection reset
			return 1, nil
		} else if strings.Contains(err.Error(), "connection refused") {
			// connection refused
			return 2, nil
		} else if strings.Contains(err.Error(), "i/o timeout") {
			// connection timeout
			return 3, nil
		} else if strings.Contains(err.Error(), "Certificate is valid for") {
			return 5, nil
		} else {
			// unknown error
			return -10, err
		}
	}

	resultCode := resp.StatusCode
	respHeader := resp.Header

	if resultCode >= 400 {
		// check 4xx - 5xx
		return resultCode, nil
	} else if resultCode >= 300 {
		// check if redirection url is correct
		redirectURL := respHeader["Location"][0]
		if redirectURL != "" {
			urlCompare, _ := url2.Parse(redirectURL)
			urlOriginElements := strings.Split(domain, ".")
			urlCompareElements := strings.Split(urlCompare.Host, ".")

			intersection := make(map[string]bool)
			for _, element := range urlOriginElements {
				intersection[element] = true
			}

			var result []string
			for _, element := range urlCompareElements {
				if intersection[element] {
					result = append(result, element)
				}
			}

			if len(result) < 2 {
				if !strings.Contains(redirectURL, domain) {
					fmt.Println("[*] BlindTLS: redirected to " + redirectURL)
					return resultCode, nil
				}
			}
		}
	}

	return 0, nil
}
