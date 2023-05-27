package http_tester

import (
	"github.com/bbalet/stopwords"
	"github.com/pywc/shawshank_intel/config"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type FilteredTokens struct {
	token       string
	resultCode  int
	redirectURL string
}

func CheckHTMLTokens(domain string) ([]FilteredTokens, error) {
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
		if len(token) < 3 || slices.Contains(testList, token) {
			continue
		}

		testList = append(testList, token)
	}

	// send each token to echo server
	filteredList := make([]FilteredTokens, 0)
	for _, token := range testList {
		req := "POST / HTTP/1.1\r\nHost: " + config.EchoServerAddr + "\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n"
		req += "magicWord=" + url.QueryEscape(token)
		resultCode, resp, redirectURL := SendHTTPRequest(config.EchoServerAddr, config.EchoServerAddr, config.EchoServerPort, req)

		if resultCode == 0 && !strings.Contains(resp, token) {
			resultCode = 399
			redirectURL = "unknown"
		} else if resultCode == 0 {
			continue
		}

		filtered := FilteredTokens{
			token:       token,
			resultCode:  resultCode,
			redirectURL: redirectURL,
		}

		filteredList = append(filteredList, filtered)
	}

	return filteredList, nil
}
