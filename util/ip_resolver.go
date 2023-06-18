package util

import (
	"errors"
	"github.com/pywc/shawshank_intel/config"
	"net"
)

func ResolveIPLocally(domain string) (string, error) {
	PrintInfo(domain, "resolving IP from Cloudflare...")

	ipList, err := net.LookupIP(domain)
	if err != nil || len(ipList) == 0 {
		newErr := errors.New("cannot resolve IP from Cloudflare")
		PrintError(config.ProxyIP, domain, newErr)
		return "", newErr
	}

	ip := ipList[0].String()

	PrintInfo(domain, "resolved IP "+ip)
	return ip, nil
}
