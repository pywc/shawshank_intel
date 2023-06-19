package util

import (
	"encoding/base64"
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

func ParseAuth() string {
	auth := base64.StdEncoding.EncodeToString([]byte(config.ProxyUsername + ":" + config.ProxyPassword))
	return auth
}

func FetchProxy(path string) ([][]string, error) {
	log.Println("Fetching proxies...")
	data, err := ReadCsvFile(path)
	if err != nil {
		PrintError(path, "", err)
		return nil, err
	}

	log.Println("Fetched " + strconv.Itoa(len(data)-1) + " proxies")

	return data[1:], nil
}
