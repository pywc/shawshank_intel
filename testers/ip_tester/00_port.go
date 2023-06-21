package ip_tester

import (
	"github.com/pywc/shawshank_intel/util"
	"net"
	"strconv"
)

func CheckTCPHandshake(ip string, port int) int {
	conn, err := util.ConnectViaProxy(ip, port, "ip")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			util.PrintError(ip, err)
		}
	}(conn)

	result := 0

	if err != nil {
		result = 1
	}

	util.PrintInfo(ip, "TCP handshake result for port "+strconv.Itoa(port)+" is "+strconv.Itoa(result))

	return result
}
