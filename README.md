# Go Web Crawler

A fast and efficient web crawler that performs deep health checks on websites. This tool crawls through internal links of a website and reports their HTTP status, content types, and any errors encountered.

## Features

- 🚀 Concurrent crawling with configurable concurrency limits
- 🔄 Rate limiting per host to prevent overwhelming servers
- 📊 Detailed status reporting for each URL
- 🔍 Content-type detection
- 🌐 Internal link detection and filtering
- ⚡ Efficient memory usage with synchronized map handling

## Installation

### Prerequisites

- Go 1.18 or higher
- Git

### Steps

1. Clone the repository:

```bash
git clone https://github.com/OsamaNagi/crawler.git
cd crawler
```

2. Install dependencies:

```bash
go mod tidy
```

3. Build the binary:

```bash
go build -o crawler
```

## Usage

### Basic Command

```bash
./crawler status <url> [maxConcurrency] [requestsPerHost] [rateInterval]
```

### Parameters

- `url`: The website URL to crawl (required)
- `maxConcurrency`: Maximum number of concurrent requests (default: 10)
- `requestsPerHost`: Maximum requests per host within the rate interval (default: 30)
- `rateInterval`: Time interval for rate limiting (default: 30s)

### Examples

1. Basic usage:

```bash
./crawler status example.com
```

2. With custom concurrency:

```bash
./crawler status example.com 20
```

3. With custom rate limiting:

```bash
./crawler status example.com 10 50 1m
```

## Output Format

The tool provides a detailed health status report:

```
Starting deep health check of example.com
This may take a while depending on the site size...
Health Status Report for example.com
=====================================
✓ https://example.com                    Status: 200 OK
✓ https://example.com/about              Status: 200 OK
✗ https://example.com/missing            Status: 404 Not Found
```

Symbols:

- ✓ : Successful response (Status code < 400)
- ✗ : Error or failed response (Status code >= 400)

## Configuration

The crawler uses sensible defaults but can be customized:

- Default max concurrency: 10 concurrent requests
- Default requests per host: 30 requests
- Default rate interval: 30 seconds

## Error Handling

The crawler handles various types of errors:

- Network errors
- Invalid URLs
- Timeout errors
- Non-HTML content types
- HTTP error status codes

## Contributing

Contributions are welcome! Please feel free to submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
