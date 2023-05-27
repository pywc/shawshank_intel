package config

// SOCKS5 Proxy Configuration
var ProxyIP string = "35.188.143.33"
var ProxyPort string = "2408"
var ProxyUsername string = ""
var ProxyPassword string = ""

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
