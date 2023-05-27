package http_tester

import (
	"bufio"
	"fmt"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"net/http"
	url2 "net/url"
	"strings"
)

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
func SendHTTPRequest(domain string, ip string, port int, req string) (int, string, string) {
	// Fetch Normally
	conn, err := util.ConnectNormally(ip, port)
	if err != nil {
		// echo server is not functional
		if ip == config.EchoServerAddr {
			return -9, "", ""
		}

		// IP is not accessible from the US
		return -2, "", ""
	}
	resp, err := util.SendHTTPTraffic(conn, req)
	conn.Close()
	if err != nil {
		return -3, "", ""
	}

	// Fetch via proxy
	conn, err = util.ConnectViaProxy(ip, port)
	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			// cannot connect to proxy
			return -1, "", ""
		} else if ip == config.EchoServerAddr {
			// residual censorship detection mode
			util.DetectResidual(domain, ip)
		} else {
			// unknown error
			fmt.Println(err.Error())
			return -10, "", ""
		}
	}

	resp, err = util.SendHTTPTraffic(conn, req)
	conn.Close()

	// check tcp errors
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			// connection reset
			return 1, "", ""
		} else if strings.Contains(err.Error(), "connection refused") {
			// connection refused
			return 2, "", ""
		} else if strings.Contains(err.Error(), "i/o timeout") {
			// connection timeout
			return 3, "", ""
		} else {
			// unknown error
			fmt.Println(err.Error())
			return -10, "", ""
		}
	}

	// parse http response
	reader := bufio.NewReader(strings.NewReader(resp))
	respObj, err := http.ReadResponse(reader, nil)
	if err != nil {
		return -10, "", ""
	}
	resultCode := respObj.StatusCode

	// check 4xx - 5xx
	if resultCode >= 300 {
		return resultCode, resp, ""
	}

	// check redirection
	redirectURL := respObj.Header["Location"][0]
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
			return respObj.StatusCode, resp, redirectURL
		}
	}

	return 0, resp, ""
}
