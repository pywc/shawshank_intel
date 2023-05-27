package http_tester

import (
	"golang.org/x/net/html"
	"net/http"
)

type FilteredTokens struct {
	token       string
	resultCode  int
	redirectURL string
}

func CheckHTMLTokens(domain string) (int, string) {
	// Fetch the HTML content from the domain
	respClean, err := http.Get("http://" + domain)
	if err != nil {
		return -3, ""
	}

	// Parse the HTML document
	doc, err := html.Parse(respClean.Body)
	respClean.Body.Close()
	if err != nil {
		return -10, ""
	}
}
