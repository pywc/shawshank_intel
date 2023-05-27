package util

import (
	"github.com/pywc/shawshank_intel/config"
	"strconv"
)

func SetProxy(ip string, port int, username string, password string) {
	config.ProxyIP = ip
	config.ProxyPort = strconv.Itoa(port)
	config.ProxyUsername = username
	config.ProxyPassword = password
}

func ParseProxyFull(remoteDNS bool) string {
	result := "socks5"

	if remoteDNS == true {
		result += "h://"
	} else {
		result += "://"
	}

	if config.ProxyUsername != "" {
		result += config.ProxyUsername
		result += ":"
		result += config.ProxyPassword
		result += "@"
	}

	result += config.ProxyIP
	result += ":"
	result += config.ProxyPort

	return result
}

func ParseProxy() string {
	result := config.ProxyIP
	result += ":"
	result += config.ProxyPort

	return result
}
