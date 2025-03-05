package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/OsamaNagi/crawler/internal/ratelimit"
)

type StatusResult struct {
	URL         string
	StatusCode  int
	ContentType string
	Error       error
}

type CrawlConfig struct {
	MaxConcurrent   int
	RequestsPerHost int
	RateInterval    time.Duration
}

func checkStatus(baseURL string, config CrawlConfig) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}

	visited := make(map[string]bool)
	results := make(chan StatusResult)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxConcurrent)
	var mu sync.Mutex

	// Add rate limiter
	rateLimiter := ratelimit.NewRateLimiter(config.RequestsPerHost, config.RateInterval)

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

		// Add rate limiting
		rateLimiter.Wait(currentURL)

		// Check status of current URL
		semaphore <- struct{}{} // Acquire
		status, contentType, err := getStatus(currentURL)
		<-semaphore // Release

		results <- StatusResult{
			URL:         currentURL,
			StatusCode:  status,
			ContentType: contentType,
			Error:       err,
		}

		// Only crawl further if it's an HTML page
		if err == nil && status < 400 && strings.Contains(contentType, "text/html") {
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
	fmt.Println("This may take a while depending on the site size...")

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

		contentInfo := ""
		if !strings.Contains(result.ContentType, "text/html") {
			contentInfo = fmt.Sprintf(" (%s)", result.ContentType)
		}

		fmt.Printf("%s %-50s Status: %d %s%s\n",
			statusSymbol,
			result.URL,
			result.StatusCode,
			statusText,
			contentInfo)
	}
}

func getStatus(url string) (int, string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, nil
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
