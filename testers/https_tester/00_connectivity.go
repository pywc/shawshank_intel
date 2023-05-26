package https_tester

import (
	"crypto/tls"
	"fmt"
	"github.com/pywc/shawshank_intel/util"
	"strings"
)

// CheckHTTPSConnectivity
/*
	-3: unhandled error
	-2: proxy error
	-1: not accessible in the US
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
*/
func CheckHTTPSConnectivity(domain string, ip string) int {
	// request configuration
	request := "GET / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"
	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
		ServerName:         domain,
	}

	// send traffic normally
	conn, err := util.ConnectNormally(ip, 443)
	if err != nil {
		return -1
	}

	_, err = util.SendHTTPSTraffic(conn, request, tlsConfig)
	conn.Close()
	if err != nil {
		return -1
	}

	// send traffic through proxy
	conn, err = util.ConnectViaProxy(ip, 443)
	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			return -2
		} else {
			return -3
		}
	}

	result, err := util.SendHTTPSTraffic(conn, request, tlsConfig)
	conn.Close()
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			return 1
		} else if strings.Contains(err.Error(), "connection reset") {
			return 2
		} else if strings.Contains(err.Error(), "i/o timeout") {
			return 3
		} else {
			return -10
		}
	}

	fmt.Println(result)

	return 0
}
