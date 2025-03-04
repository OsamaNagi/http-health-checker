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

	visited := make(map[string]bool)
	results := make(chan StatusResult)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrent)
	var mu sync.Mutex

	// Recursive function to check links
	var checkLinksRecursive func(currentURL string)
	checkLinksRecursive = func(currentURL string) {
		defer wg.Done()

		mu.Lock()
		if visited[currentURL] {
			mu.Unlock()
			return
		}
		visited[currentURL] = true
		mu.Unlock()

		// Check status of current URL
		semaphore <- struct{}{} // Acquire
		status, err := getStatus(currentURL)
		<-semaphore // Release

		results <- StatusResult{
			URL:        currentURL,
			StatusCode: status,
			Error:      err,
		}

		// If this page is accessible, get its links
		if err == nil && status < 400 {
			links, err := getAllLinks(currentURL)
			if err != nil {
				fmt.Printf("Error getting links from %s: %v\n", currentURL, err)
				return
			}

			// Process internal links
			for _, link := range links {
				if !isInternalLink(parsedURL, link) {
					continue
				}

				mu.Lock()
				notVisited := !visited[link]
				mu.Unlock()

				if notVisited {
					wg.Add(1)
					go checkLinksRecursive(link)
				}
			}
		}
	}

	// Start the recursive crawl
	fmt.Printf("\nStarting deep health check of %s\n", baseURL)
	fmt.Println("This may take a while depending on the site size...\n")

	wg.Add(1)
	go checkLinksRecursive(baseURL)

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	fmt.Printf("Health Status Report for %s\n", baseURL)
	fmt.Println("=====================================")

	for result := range results {
		if result.Error != nil {
			fmt.Printf("%-50s Error: %v\n", result.URL, result.Error)
			continue
		}
		statusText := http.StatusText(result.StatusCode)
		statusSymbol := "✓"
		if result.StatusCode >= 400 {
			statusSymbol = "✗"
		}
		fmt.Printf("%s %-50s Status: %d %s\n", statusSymbol, result.URL, result.StatusCode, statusText)
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
