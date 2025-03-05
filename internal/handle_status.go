package internal

import (
	"fmt"
	"strconv"
	"time"
)

const (
	defaultMaxConcurrency  = 10
	defaultRequestsPerHost = 30
	defaultRateInterval    = 30 * time.Second
)

func HandleStatus(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: URL is required for status check")
		return
	}

	url := args[0]
	config := CrawlConfig{
		MaxConcurrent:   defaultMaxConcurrency,
		RequestsPerHost: defaultRequestsPerHost,
		RateInterval:    defaultRateInterval,
	}

	if len(args) >= 2 {
		if mc, err := strconv.Atoi(args[1]); err == nil {
			config.MaxConcurrent = mc
		}
	}
	if len(args) >= 3 {
		if rph, err := strconv.Atoi(args[2]); err == nil {
			config.RequestsPerHost = rph
		}
	}
	if len(args) >= 4 {
		if d, err := time.ParseDuration(args[3]); err == nil {
			config.RateInterval = d
		}
	}

	CheckStatus(url, config)
}

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  crawler status <url> [maxConcurrency] [requestsPerHost] [rateInterval]")
	fmt.Printf("Defaults: maxConcurrency=%d, requestsPerHost=%d, rateInterval=%v\n",
		defaultMaxConcurrency,
		defaultRequestsPerHost,
		defaultRateInterval)
}
