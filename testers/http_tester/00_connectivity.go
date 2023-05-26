package http_tester

import (
	"fmt"
	"github.com/pywc/shawshank_intel/util"
	"strings"
)

// CheckHTTPConnectivity
/*
	-2: unhandled error
	-1: proxy error
	0: success
	1: reset
	2: refused
	3: silent drop
	4: TODO: throttle
	5: redirection
*/
func CheckHTTPConnectivity(domain string) int {
	request := "GET / HTTP/1.1\r\nHost: " + domain + "\r\n\r\n"

	conn, err := util.ConnectViaProxy(domain, 80)
	if err != nil {
		if strings.Contains(err.Error(), "general SOCKS server failure") {
			return -1
		} else {
			return -2
		}
	}

	result, err := util.SendHTTPTraffic(conn, request)
	conn.Close()

	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			return 3
		} else {
			return -10
		}
	}

	fmt.Println(result)

	conn.Close()

	return 0
}
