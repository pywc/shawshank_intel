package util

import (
	"github.com/pywc/shawshank_intel/config"
	"log"
	"strconv"
)

func SetProxy(ip string, port int, username string, password string, proxyType string) {
	config.ProxyIP = ip
	config.ProxyPort = strconv.Itoa(port)
	config.ProxyUsername = username
	config.ProxyPassword = password
	config.ProxyType = proxyType

	log.Println("Current proxy is set to " + ip + " of type " + proxyType)
}

func ParseSOCKS5ProxyFull(remoteDNS bool) string {
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

func FetchProxy(path string) ([][]string, error) {
	log.Println("Fetching proxies...")
	data, err := ReadCsvFile(path)
	if err != nil {
		PrintError(path, "", err)
		return nil, err
	}

	log.Println("Fetched " + string(len(data)) + " proxies")

	return data, nil
}
