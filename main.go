package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultMaxConcurrency = 10
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "status":
		handleStatus(os.Args[2:])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  crawler status <url> [maxConcurrency]")
	fmt.Printf("Defaults: maxConcurrency=%d\n",
		defaultMaxConcurrency)
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
