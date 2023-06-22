package ip_tester

import (
	"github.com/pywc/shawshank_intel/util"
	"strconv"
)

func CheckTCPHandshake(ip string, port int) int {
	result := 0

	conn, err := util.ConnectViaProxy(ip, port, "ip")
	if err != nil {
		result = 1
	} else {
		conn.Close()
	}

	util.PrintInfo(ip, "TCP handshake result for port "+strconv.Itoa(port)+" is "+strconv.Itoa(result))

	return result
}
