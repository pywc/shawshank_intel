package http_tester

import (
	"fmt"
	"github.com/jpillora/go-tld"
	combinations "github.com/mxschmitt/golang-combinations"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"golang.org/x/exp/slices"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type RequestWord struct {
	Hostname          string
	GetWord           string `default:"GET"`
	HttpWord          string `default:"HTTP/1.1"`
	HostWord          string `default:"Host:"`
	HttpDelimiterWord string `default:"\r\n"`
	Path              string `default:"/"`
	Header            string `default:""`
}

func FormatHttpRequest(requestWord RequestWord) string {
	getWord := "GET"
	if requestWord.GetWord != "" {
		getWord = requestWord.GetWord
	}

	httpWord := "HTTP/1.1"
	if requestWord.HttpWord != "" {
		httpWord = requestWord.HttpWord
	}

	hostWord := "Host: "
	if requestWord.HostWord != "" {
		hostWord = requestWord.HostWord
	}

	httpDelimiterWord := "\r\n"
	if requestWord.HttpDelimiterWord != "" {
		httpDelimiterWord = requestWord.HttpDelimiterWord
	}

	path := ""
	if config.ProxyType == "https" {
		path = " http://" + requestWord.Hostname + "/ "
		if requestWord.Path != "" {
			path = " http://" + requestWord.Hostname + "/"
			path += strings.TrimSpace(requestWord.Path)
		}
	} else {
		path = " / "
		if requestWord.Path != "" {
			path = requestWord.Path
		}
	}

	header := ""
	if requestWord.Header != "" {
		header = requestWord.Header + "\r\n"
	}

	//Handle hostname changes - This has to be done at runtime, since the strategies would be selected first, but the hostname itself is only known at runtime
	var host string
	hostNameParts := strings.Split(requestWord.Hostname, "|")
	if len(hostNameParts) > 1 {
		//ServerNameParts[1] contains the strategy to be run at runtime
		if hostNameParts[1] == "omit" {
			format := "%s%s%s%s\r\n%s\r\n"
			return fmt.Sprintf(format, getWord, path, httpWord, httpDelimiterWord, header)
		} else if hostNameParts[1] == "empty" {
			host = ""
		} else if hostNameParts[1] == "repeat" {
			//Now there should be a third part that says how many times to repeat
			repeatTimes, err := strconv.Atoi(hostNameParts[2])
			if err != nil {
				log.Println("[http_fuzzer.FormatHTTPRequest] Error converting string into integer (repeat)")
				log.Println(err)
				log.Println("Reverting to default")
				host = hostNameParts[0]
			} else {
				host = util.Repeat(hostNameParts[0], repeatTimes)
			}
		} else if hostNameParts[1] == "reverse" {
			host = util.Reverse(hostNameParts[0])
		}
	} else {
		host = requestWord.Hostname
	}

	format := "%s%s%s%s%s%s\r\n%s"
	if config.ProxyType == "https" && config.ProxyUsername != "" {
		format += "Proxy-Authorization: Basic " + util.ParseAuth() + "\r\n"
	}
	format += "\r\n"

	result := fmt.Sprintf(format, getWord, path, httpWord, httpDelimiterWord, hostWord, host, header)
	return result
}

func GenerateRandomCapitalizedValues(word string) string {
	randomlyCapitalizedWord := ""
	for _, char := range word {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			randomlyCapitalizedWord += string(char)
			continue
		}
		rand.Seed(time.Now().UTC().UnixNano())
		choice := rand.Intn(2)
		if choice == 0 {
			randomlyCapitalizedWord += strings.ToLower(string(char))
		} else if choice == 1 {
			randomlyCapitalizedWord += strings.ToUpper(string(char))
		}
	}
	return randomlyCapitalizedWord
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CapitalizedPermutations(ip string, op string) []string {
	var s []string
	if len(ip) == 0 {
		return []string{op}
	}
	lowerChar := strings.ToLower(string(ip[0]))
	upperChar := strings.ToUpper(string(ip[0]))
	ip = ip[1:len(ip)]
	s = append(s, CapitalizedPermutations(ip, op+lowerChar)...)
	s = append(s, CapitalizedPermutations(ip, op+upperChar)...)
	return unique(s)
}

func GenerateAllCapitalizedPermutations(word string) []string {
	return CapitalizedPermutations(word, "")
}

func GenerateRandomlyRemovedWord(word string) string {
	randomlyRemovedWord := ""
	for _, char := range word {
		rand.Seed(time.Now().UTC().UnixNano())
		choice := rand.Intn(2)
		if choice == 1 {
			randomlyRemovedWord += string(char)
		}
	}
	return randomlyRemovedWord
}

func GenerateAllSubstringPermutations(word string) []string {
	splitWord := strings.Split(word, "")
	combs := combinations.All(splitWord)
	var permutations []string
	for _, elem := range combs {
		permutations = append(permutations, strings.Join(elem, ""))
	}
	return permutations
}

func GenerateAlternatives(alternatives []string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	choice := rand.Intn(len(alternatives))
	return alternatives[choice]
}

func GenerateAllAlternatives(alternatives []string) []string {
	return alternatives
}

func GenerateHostNameRandomPadding() string {
	prefixPaddingLength := rand.Intn(5)
	suffixPaddingLength := rand.Intn(5)
	hostnameWithPadding := strings.Repeat("*", prefixPaddingLength)
	hostnameWithPadding += "%s"
	hostnameWithPadding += strings.Repeat("*", suffixPaddingLength)
	return hostnameWithPadding
}

func GenerateAllHostNamePaddings(hostname string) []string {
	var hostnameWithAllPadding []string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			hostnameWithPadding := strings.Repeat("*", i)
			hostnameWithPadding += hostname
			hostnameWithPadding += strings.Repeat("*", j)
			hostnameWithAllPadding = append(hostnameWithAllPadding, hostnameWithPadding)
		}

	}
	return hostnameWithAllPadding
}

var GetAlternatives = []string{"POST", "PUT", "PATCH", "DELETE", "XXX", " "}

func GenerateGetAlternatives() string {
	return GenerateAlternatives(GetAlternatives)
}

func GenerateAllGetAlternatives() []string {
	allList := GenerateAllAlternatives(GetAlternatives)

	if config.ProxyType == "https" {
		noHTTPSProxy := []string{" ", "XXX"}
		newList := make([]string, 0)

		for _, entry := range allList {
			if slices.Contains(noHTTPSProxy, entry) {
				continue
			}

			newList = append(newList, entry)
		}

		return newList
	}

	return allList
}

var HttpAlternatives = []string{"XXXX/1.1", "HTTP/11.1", "HTTP/1.12", "/11.1", "HTTP2", "HTTP3", "HTTP9", "HTTP/2", "HTTP/3", "HTTP/9", " ", "HTTPx/1.1", "HTTP /1.1", "HTTP/ 1.1", "HTTP/1.1x", "HTTP/x1.1"}

func GenerateHttpAlternatives() string {
	return GenerateAlternatives(HttpAlternatives)
}

func GenerateAllHttpAlternatives() []string {
	return GenerateAllAlternatives(HttpAlternatives)
}

var HostAlternatives = []string{"XXXX: ", "XXXX:", "Host:\r\n", "Hostwww.", "Host:www.", "HostHeader:", " "}

func GenerateHostAlternatives() string {
	return GenerateAlternatives(HostAlternatives)
}

func GenerateAllHostAlternatives() []string {
	return GenerateAllAlternatives(HostAlternatives)
}

var PathAlternatives = []string{"/ ", " z ", " ? ", " ", " /", "**", " /x", "x/ "}

func GeneratePathAlternatives() string {
	return GenerateAlternatives(PathAlternatives)
}

func GenerateAllPathAlternatives(hostname string) []string {
	return GenerateAllAlternatives(PathAlternatives)
}

var HTTPHeaders = []string{"Accept: text/html", "Accept: application/xml", "Accept: text/html,application/xhtml+xml", "Accept: application/json", "Accept: xxx", "Accept-Charset: utf-8", "Accept-Charset: xxx", "Accept-Datetime: Thu, 31 May 2007 20:35:00 GMT", "Accept-Datetime: xxx", "Accept-Encoding: gzip, deflate", "Accept-Encoding: xxx", "Accept-Language: en-US", "Accept-Language: xxx", "Access-Control-Request-Method: GET", "Access-Control-Request-Method: xxx", "Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "Cache-Control: no-cache", "Cache-Control: xxx", "Connection: keep-alive", "Connection: xxx", "Content-Encoding: gzip", "Content-Encoding: xxx", "Content-Length: 1000", "Content-MD5: Q2hlY2sgSW50ZWdyaXR5IQ==", "Content-Type: application/x-www-form-urlencoded", "Content-Type: xxx", "Cookie: $Version=1; Skin=new;", "Cookie: xxx", "Date: Tue, 15 Nov 1994 08:12:31 GMT", "Expect: 100-continue", "Expect: xxx", "From: user@example.com", "If-Match: \"737060cd8c284d8af7ad3082f209582d\"", "If-Modified-Since: Sat, 29 Oct 1994 19:43:31 GMT", "If-None-Match: \"737060cd8c284d8af7ad3082f209582d]\"", "If-Range: \"737060cd8c284d8af7ad3082f209582d\"", "If-Unmodified-Since: Sat, 29 Oct 1994 19:43:31 GMT", "Max-Forwards: 10", "Max-Forwards: xxx", "Origin: http://www.example-xxx.com", "Pragma: no-cache", "Pragma: xxx", "Prefer: return=representation", "Prefer: xxx", "Proxy-Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==", "Range: bytes=500-999", "Referer: http://example-xxx.com", "TE: trailers, deflate", "Trailer: Max-Forwards", "Trailer: xxx", "Transfer-Encoding: chunked", "Transfer-Encoding: xxx", "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:12.0) Gecko/20100101 Firefox/12.0", "User-Agent: xxx", "Upgrade: h2c, HTTPS/1.3, IRC/6.9, RTA/x11, websocket", "Upgrade: xxx", "Via: 1.0 fred, 1.1 example-xxx.com (Apache/1.1)", "Warning: 199 Miscellaneous warning", "Warning: xxx"}

func GenerateHeaderAlternatives() string {
	return GenerateAlternatives(HTTPHeaders)
}

func GenerateAllHeaderAlternatives() []string {
	return GenerateAllAlternatives(HTTPHeaders)
}

// https://azbigmedia.com/business/here-are-2021s-most-popular-tlds-and-domain-registration-trends/
var TLDs = []string{"com", "xyz", "net", "club", "me", "org", "co", "shop", "info", "live"}

// https://securitytrails.com/blog/most-popular-subdomains-mx-records#:~:text=As%20you%20can%20see%2C%20the,forums%2C%20wiki%2C%20community).
var Subdomains = []string{"www", "mail", "forum", "m", "blog", "shop", "forums", "wiki", "community", "ww1"}

var hostnames = []string{"%s|omit", "%s|empty", "%s|repeat|2", "%s|repeat|3", "%s|reverse"}

func GenerateTLDAlternatives() string {
	return GenerateAlternatives(TLDs)
}

func GenerateAllTLDAlternatives(hostname string) []string {
	tldAlternatives := GenerateAllAlternatives(TLDs)

	for i, alt := range tldAlternatives {
		u, _ := tld.Parse("https://" + hostname)
		newUrl := ""
		if u.Subdomain != "" {
			newUrl += u.Subdomain + "."
		}

		newUrl += u.Domain + "." + alt
		tldAlternatives[i] = newUrl
	}

	return tldAlternatives
}

func GenerateSubdomainsAlternatives() string {
	return GenerateAlternatives(Subdomains)
}

func GenerateAllSubdomainsAlternatives(hostname string) []string {
	subdomainAlternatives := GenerateAllAlternatives(Subdomains)

	for i, alt := range subdomainAlternatives {
		u, _ := tld.Parse("https://" + hostname)
		newUrl := alt + "." + u.Domain
		if u.TLD != "" {
			newUrl += "." + u.TLD
		}

		subdomainAlternatives[i] = newUrl
	}

	return subdomainAlternatives
}

func GenerateHostNameAlternatives() string {
	return GenerateAlternatives(hostnames)
}

func GenerateAllHostNameAlternatives(hostname string) []string {
	hostnameAlternatives := GenerateAllAlternatives(hostnames)

	for i, alt := range hostnameAlternatives {
		hostnameAlternatives[i] = fmt.Sprintf(alt, hostname)
	}

	return hostnameAlternatives
}
