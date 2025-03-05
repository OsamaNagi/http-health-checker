package ratelimit

import (
	"net/url"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	domains  map[string]*domainLimiter
	perHost  int
	interval time.Duration
}

type domainLimiter struct {
	tokens    int
	lastReset time.Time
}

func NewRateLimiter(requestsPerHost int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		domains:  make(map[string]*domainLimiter),
		perHost:  requestsPerHost,
		interval: interval,
	}
}

func (rl *RateLimiter) Wait(urlStr string) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	host := parsedURL.Host

	rl.mu.Lock()
	limiter, exists := rl.domains[host]
	if !exists {
		limiter = &domainLimiter{
			tokens:    rl.perHost,
			lastReset: time.Now(),
		}
		rl.domains[host] = limiter
	}

	for {
		now := time.Now()
		elapsed := now.Sub(limiter.lastReset)

		if elapsed >= rl.interval {
			limiter.tokens = rl.perHost
			limiter.lastReset = now
		}

		if limiter.tokens > 0 {
			limiter.tokens--
			rl.mu.Unlock()
			return
		}

		rl.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		rl.mu.Lock()
	}
}
