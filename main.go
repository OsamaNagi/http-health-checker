package main

import (
	"os"

	. "github.com/OsamaNagi/http-health-checker/internal"
)

func main() {
	if len(os.Args) < 2 {
		PrintUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "status":
		HandleStatus(os.Args[2:])
	case "help":
		PrintUsage()
	}
}
