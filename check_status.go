package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type StatusResult struct {
	URL        string
	StatusCode int
	Error      error
}

func checkStatus(baseURL string, maxConcurrent int) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}

	// Get all links first
	urls, err := getAllLinks(baseURL)
	if err != nil {
		fmt.Printf("Error getting links: %v\n", err)
		return
	}

	results := make(chan StatusResult)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrent)

	// Start workers for each URL
	for _, link := range urls {
		if !isInternalLink(parsedURL, link) {
			continue
		}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			status, err := getStatus(url)
			results <- StatusResult{
				URL:        url,
				StatusCode: status,
				Error:      err,
			}
		}(link)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	fmt.Printf("\nHealth Status Report for %s\n", baseURL)
	fmt.Println("=====================================")

	for result := range results {
		if result.Error != nil {
			fmt.Printf("%-50s Error: %v\n", result.URL, result.Error)
			continue
		}
		statusText := http.StatusText(result.StatusCode)
		fmt.Printf("%-50s Status: %d %s\n", result.URL, result.StatusCode, statusText)
	}
}

func getStatus(url string) (int, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func getAllLinks(baseURL string) ([]string, error) {
	htmlBody, err := getHTML(baseURL)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return getURLsFromHTML(htmlBody, parsedURL)
}

func isInternalLink(baseURL *url.URL, link string) bool {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return false
	}
	return parsedLink.Hostname() == baseURL.Hostname()
}
