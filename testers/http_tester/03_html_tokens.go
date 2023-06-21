package http_tester

import (
	"github.com/bbalet/stopwords"
	"github.com/pywc/shawshank_intel/config"
	"github.com/pywc/shawshank_intel/util"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CheckHTMLTokens(domain string) ([]FilteredHTTP, error) {
	// Fetch the HTML content from the domain
	respClean, err := http.Get("http://" + domain)
	if err != nil {
		return nil, err
	}

	// Parse the HTML document
	doc, err := io.ReadAll(respClean.Body)
	respClean.Body.Close()
	if err != nil {
		return nil, err
	}

	// tokenize html document and remove stopwords
	tempStr := stopwords.CleanString(string(doc), "en", true)
	tempList := strings.Split(tempStr, " ")
	testList := make([]string, 0)
	for _, token := range tempList {
		token = config.NonAlphanumericRegex.ReplaceAllString(token, "")
		if len(token) < 4 || slices.Contains(testList, token) {
			continue
		}

		testList = append(testList, token)
	}

	// send each token to echo server
	filteredList := make([]FilteredHTTP, 0)
	for _, token := range testList {
		reqBody := "magicWord=" + url.QueryEscape(token)
		req := "POST http://" + util.ParseEcho() + " HTTP/1.1\r\n" +
			"Host: " + util.ParseEcho() + "\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: " + strconv.Itoa(len(reqBody)) + "\r\n"
		if config.ProxyUsername != "" {
			req += "Proxy-Authorization: Basic " + util.ParseAuth() + "\r\n"
		}
		req += "\r\n"
		req += reqBody
		resultCode, resp, redirectURL, err := SendHTTPRequest(config.EchoServerAddr, config.EchoServerAddr, config.EchoServerPort, req)
		if resultCode == -10 {
			log.Println("[*] Error - " + domain + " - " + err.Error())
		}

		if resultCode == 0 && !strings.Contains(resp, url.QueryEscape(token)) {
			resultCode = 399
			redirectURL = "unknown"
		} else if resultCode == 0 {
			util.PrintInfo(domain, "html token result for \""+token+"\": "+strconv.Itoa(resultCode))
			continue
		}

		filtered := FilteredHTTP{
			Component:   token,
			ResultCode:  resultCode,
			RedirectURL: redirectURL,
		}

		util.PrintInfo(domain, "html token result for \""+token+"\": "+strconv.Itoa(resultCode))

		filteredList = append(filteredList, filtered)
	}

	return filteredList, nil
}
