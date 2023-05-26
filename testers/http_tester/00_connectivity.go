package http_tester

import (
	"github.com/pywc/shawshank_intel/util"
)

// CheckHTTPConnectivity
/*
	-1: not accessible in the US
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
	result, err := util.SendHTTPTraffic(conn, request)
	conn.Close()

	return 0
}
