package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
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

func main() {
	// Fetch the HTML content from a URL
	resp, err := http.Get("http://naver.com") // Replace with the desired URL
	if err != nil {
		log.Fatal("Failed to fetch HTML content:", err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse HTML:", err)
	}

	// Find the title element
	title := findTitle(doc)
	fmt.Println(title)
}
