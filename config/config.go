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

// Residual Test Configuration
var ResidualTestThreshold float64 = 600

// HTTP Test Configuration
var EchoServerAddr string = "35.188.143.33"
var EchoServerPort int = 8008
var MagicWord string = "somethingsSpeciale"

// HTTPS Test Configuration
var DummyServerDomain = "wisc.edu"
var DummyServerIP = "144.92.9.70"

// For both HTTP and HTTPS
func DomainWildcards(domain string) []string {
	testList := make([]string, 7)
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
