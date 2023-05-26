package ip

import (
	"github.com/pywc/shawshank_intel/util"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"strconv"
)

func TestPort(ip string, port int) int {
	dialer, err := proxy.SOCKS5("tcp", util.ParseProxy(), nil, proxy.Direct)
	if err != nil {
		log.Fatal("Failed to create proxy dialer:", err)
		return -1
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		log.Fatal("Failed to connect to the destination:", err)
		return 1
	}
	defer conn.Close()

	return 0
}
