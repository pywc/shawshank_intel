package config

import (
	"regexp"
	"time"
)

// Global Configuration
var TestCount = 30
var CurrentComponent = ""
var Timeout = time.Second * 20

// SOCKS5 Proxy Configuration
var ProxyIP = "35.188.143.33"
var ProxyPort = "2408"
var ProxyUsername = ""
var ProxyPassword = ""
var ProxyType = "socks5"
var ProxyCountry = "us"
var ProxyISP = "charter"

// Throttle Test Configuration
var ThrottlePValThreshold = 0.01
var TDigestCompression float64 = 1000

// Residual Test Configuration
var ResidualTestThreshold float64 = 600

// HTTP Test Configuration
var EchoServerAddr string = "35.188.143.33"
var EchoServerPort int = 80
var MagicWord string = "somethingSpecial"
var DomainSimilarityThreshold float64 = 0.5
var NonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// HTTPS Test Configuration
// not wisc.edu since it returns handshake failure for some reason
var DummyServerDomain = "mit.edu"
var DummyServerIP = "23.64.108.89"
var ESNIDomain = "www.cloudflare.com"
var ESNIIP = "104.16.123.96"

// For both HTTP and HTTPS
func DomainWildcards(domain string) []string {
	testList := make([]string, 0)
	testList = append(testList, domain)
	testList = append(testList, "abc123qwezxc."+domain)
	testList = append(testList, "*."+domain)
	testList = append(testList, "*"+domain)
	testList = append(testList, domain+".com")
	testList = append(testList, domain+".*")
	testList = append(testList, domain+"*")

	return testList
}

var UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/114.0"
