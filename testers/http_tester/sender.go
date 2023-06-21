package http_tester

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/jpillora/go-tld"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"io"
	"strconv"
	"strings"
)

type FilteredHTTP struct {
	Component   string `json:"component,omitempty"`
	ResultCode  int    `json:"result_code,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
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
		util.PrintError(domain, err)
		return -2, "", "", err
	}

	reqNormal := "GET / HTTP/1.1\r\n" +
		"Host: " + domain + "\r\n" +
		"Accept: */*\r\n" +
		"User-Agent: " + config.UserAgent + "\r\n\r\n"

	resp, err := util.SendHTTPTraffic(conn, reqNormal)
	conn.Close()
	if err != nil {
		util.PrintError(domain, err)
		return -3, "", "", err
	}

	// Fetch via proxy
	conn, err = util.ConnectViaProxy(ip, port, "http")
	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			// cannot connect to proxy
			return -1, "", "", nil
		} else if ip == config.EchoServerAddr {
			// residual censorship detection mode
			util.DetectResidual(ip, port, "http")
		} else {
			// unknown error
			util.PrintError(domain, err)
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
			return -10, "", "", err
		}
	}

	resultCode := resp.StatusCode
	respHeader := resp.Header
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	conn.Close()

	if resultCode == 521 {
		// cloudflare server down
		return -3, string(respBody), "", nil
	} else if resultCode >= 400 {
		// check 4xx - 5xx
		return resultCode, string(respBody), "", nil
	} else if resultCode >= 300 {
		// check redirection
		redirectURL := respHeader["Location"]
		if len(redirectURL) < 1 {
			return resultCode, string(respBody), "unknown", nil
		} else if redirectURL[0] != "" {
			urlCompare, err := tld.Parse(redirectURL[0])
			if err != nil {
				util.PrintError(domain, err)
				return resultCode, string(respBody), "unknown", nil
			}

			urlOriginal, err := tld.Parse("http://" + domain)
			if err != nil {
				util.PrintError(domain, err)
				return resultCode, string(respBody), "unknown", nil
			}

			similarity := strutil.Similarity(urlOriginal.Domain, urlCompare.Domain, metrics.NewHamming())
			util.PrintInfo(domain, "redirection: similarity between \""+
				urlOriginal.Domain+"\" and \""+urlCompare.Domain+"\" is "+
				strconv.FormatFloat(similarity, 'f', -1, 64))

			if similarity < config.DomainSimilarityThreshold {
				return resultCode, string(respBody), redirectURL[0], nil
			}
		}
	}

	return 0, string(respBody), "", nil
}
