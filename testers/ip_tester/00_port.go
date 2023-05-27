package ip_tester

import (
	"github.com/pywc/shawshank_intel/util"
	"log"
)

func TestPort(ip string, port int) int {
	conn, err := util.ConnectViaProxy(ip, port)
	defer conn.Close()

	if err != nil {
		log.Fatal("Failed to connect to the destination:", err)
		return 1
	}

	return 0
}
