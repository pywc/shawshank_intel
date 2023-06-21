package http_tester

import (
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Recursively find the title element in the HTML document
func findTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" {
		if node.FirstChild != nil {
			return node.FirstChild.Data
		}
	}

	// Traverse child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result := findTitle(child)
		if result != "" {
			return result
		}
	}

	return ""
}

// CheckHTMLTitle Checks whether HTML <title> tag is used for filtering
// Sets 399 as result code if it does not contain the correct magic word
func CheckHTMLTitle(domain string) []FilteredHTTP {
	// Fetch the HTML content from the domain
	respClean, err := http.Get("http://" + domain)
	if err != nil {
		return nil
	}

	// Parse the HTML document
	doc, err := html.Parse(respClean.Body)
	respClean.Body.Close()
	if err != nil {
		return nil
	}

	// Find the title element
	title := findTitle(doc)

	// Send the title under <title> to echo server
	html := "<html><head><title>" + title + "</title></head><body>" + config.MagicWord + "</body></html>"
	reqBody := "magicWord=" + url.QueryEscape(html)
	req := "POST http://" + config.EchoServerAddr + " HTTP/1.1\r\n" +
		"Host: " + util.ParseProxy() + "\r\n" +
		"Content-Type: application/x-www-form-urlencoded\r\n" +
		"Content-Length: " + strconv.Itoa(len(reqBody)) + "\r\n"
	if config.ProxyUsername != "" {
		req += "Proxy-Authorization: Basic " + util.ParseAuth() + "\r\n"
	}
	req += "\r\n"
	req += reqBody
	resultCode, respEcho, redirectURL, err := SendHTTPRequest(config.EchoServerAddr,
		config.EchoServerAddr, config.EchoServerPort, req, "")
	if resultCode == -10 {
		log.Println("[*] Error - " + domain + " - " + err.Error())
	}

	if resultCode == 0 && !strings.Contains(respEcho, config.MagicWord) {
		resultCode = 399
		redirectURL = "unknown"
	}

	util.PrintInfo(domain, "html title result for \""+title+"\" is "+strconv.Itoa(resultCode))

	return []FilteredHTTP{
		{
			Component:   title,
			ResultCode:  resultCode,
			RedirectURL: redirectURL,
		},
	}
}
