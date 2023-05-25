package util

import (
	"../config"
	"strconv"
)

func parse_proxy(remoteDNS bool) string {
	var result string = "socks5"

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
