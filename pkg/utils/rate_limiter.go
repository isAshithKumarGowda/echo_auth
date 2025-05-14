package utils

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	maxRequest = 5
	timeWidow  = 24 * time.Hour
)

var (
	mu       sync.Mutex
	limiters = make(map[string]*rate.Limiter)
)

func GetRateLimiter(email string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := limiters[email]; !exists {
		limiters[email] = rate.NewLimiter(rate.Every(timeWidow), maxRequest)
	}

	return limiters[email]
}
