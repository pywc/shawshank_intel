package util

import (
	"../config"
	"strconv"
)

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
	result += strconv.Itoa(config.ProxyPort)

	return result
}

func ParseProxy() string {
	result := config.ProxyIP
	result += ":"
	result += strconv.Itoa(config.ProxyPort)

	return result
}
