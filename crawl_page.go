package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse URL: %v\n", err)
		return
	}

	rawURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("couldn't parse URL: %v\n", err)
		return
	}

	if currentURL.Host != rawURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't normalize URL: %v\n", err)
		return
	}

	if _, visited := pages[normalizedURL]; visited {
		pages[normalizedURL]++
		return
	}

	pages[normalizedURL] = 1

	fmt.Printf("crawling %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v", err)
		return
	}

	nextURLs, err := getURLsFromHTML(htmlBody, rawBaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		crawlPage(rawBaseURL, nextURL, pages)
	}
}
