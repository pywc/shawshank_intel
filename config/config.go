package config

// Global Configuration
var TestCount = 30

// SOCKS5 Proxy Configuration
var ProxyIP = "35.188.143.33"
var ProxyPort = "2408"
var ProxyUsername = ""
var ProxyPassword = ""

// Throttle Test Configuration
var ThrottlePValThreshold = 0.01
var TDigestCompression float64 = 1000

// HTTP Test Configuration
var EchoServerAddr string = "35.188.143.33"
var EchoServerPort int = 8008
var MagicWord string = "somethingsSpeciale"

func HTTPHostWildcards(domain string) []string {
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
