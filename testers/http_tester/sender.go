package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"io"
	url2 "net/url"
	"strings"
)

type FilteredHTTP struct {
	component   string
	resultCode  int
	redirectURL string
}

// SendHTTPRequest Returns result_code, response_body, redirect_url (if redirection)
/*
	Result Code Entry
	=========================
	-10: unhandled error
	-9: error in echo server
	-3: Server not open to US
	-2: IP not open to US
	-1: proxy error
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
	3xx: redirection
	4xx: not accessible
	5xx: internal server error
*/
func SendHTTPRequest(domain string, ip string, port int, req string) (int, string, string, error) {
	// Fetch Normally
	conn, err := util.ConnectNormally(ip, port)
	if err != nil {
		// echo server is not functional
		if ip == config.EchoServerAddr {
			return -9, "", "", nil
		}

		// IP is not accessible from the US
		return -2, "", "", err
	}

	resp, err := util.SendHTTPTraffic(conn, req)
	conn.Close()
	if err != nil {
		return -3, "", "", err
	}

	// Fetch via proxy
	conn, err = util.ConnectViaProxy(ip, port)
	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			// cannot connect to proxy
			return -1, "", "", nil
		} else if ip == config.EchoServerAddr {
			// residual censorship detection mode
			util.DetectResidual(domain, ip, port)
		} else {
			// unknown error
			return -10, "", "", err
		}
	}

	resp, err = util.SendHTTPTraffic(conn, req)

	// check tcp errors
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			// connection reset
			return 1, "", "", nil
		} else if strings.Contains(err.Error(), "connection refused") {
			// connection refused
			return 2, "", "", nil
		} else if strings.Contains(err.Error(), "i/o timeout") {
			// connection timeout
			return 3, "", "", nil
		} else {
			// unknown error
			return -10, "", "", nil
		}
	}

	resultCode := resp.StatusCode
	respHeader := resp.Header
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	conn.Close()

	if resultCode >= 400 {
		// check 4xx - 5xx
		return resultCode, string(respBody), "", nil
	} else if resultCode >= 300 {
		// check redirection
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
				return resultCode, string(respBody), redirectURL, nil
			}
		}
	}

	return 0, string(respBody), "", nil
}
