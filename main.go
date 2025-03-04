package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultMaxConcurrency = 10
	defaultMaxPages       = 20
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "crawl":
		handleCrawl(os.Args[2:])
	case "status":
		handleStatus(os.Args[2:])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  crawler crawl <url> [maxConcurrency] [maxPages]")
	fmt.Println("  crawler status <url> [maxConcurrency]")
	fmt.Printf("Defaults: maxConcurrency=%d, maxPages=%d\n",
		defaultMaxConcurrency, defaultMaxPages)
}

func handleStatus(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: URL is required for status check")
		return
	}

	url := args[0]
	maxConcurrent := defaultMaxConcurrency

	if len(args) >= 2 {
		if mc, err := strconv.Atoi(args[1]); err == nil {
			maxConcurrent = mc
		}
	}

	checkStatus(url, maxConcurrent)
}

func handleCrawl(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: URL is required")
		fmt.Println("Usage: crawler <url> [maxConcurrency] [maxPages]")
		fmt.Printf("Default maxConcurrency: %d, Default maxPages: %d\n",
			defaultMaxConcurrency, defaultMaxPages)
		return
	}

	rawBaseURL := args[0]
	maxConcurrency := defaultMaxConcurrency
	maxPages := defaultMaxPages

	if len(args) >= 2 {
		if mc, err := strconv.Atoi(args[1]); err == nil {
			maxConcurrency = mc
		} else {
			fmt.Printf("Warning: invalid maxConcurrency, using default: %d\n", defaultMaxConcurrency)
		}
	}

	if len(args) >= 3 {
		if mp, err := strconv.Atoi(args[2]); err == nil {
			maxPages = mp
		} else {
			fmt.Printf("Warning: invalid maxPages, using default: %d\n", defaultMaxPages)
		}
	}

	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v\n", err)
		return
	}

	fmt.Printf("Starting crawl of: %s\n", rawBaseURL)
	fmt.Printf("Max concurrent requests: %d\n", maxConcurrency)
	fmt.Printf("Max pages to crawl: %d\n", maxPages)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	printReport(cfg.pages, rawBaseURL)
}
