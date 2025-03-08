// middleware/rate_limit.go
package middleware
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// visitor holds the rate limiter and last seen time for each visitor
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// visitors keeps track of all visitors and their rate limiters
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Initialize a background goroutine to clean up old entries
func init() {
	go cleanupVisitors()
}

// cleanupVisitors removes old entries from the visitors map
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// getVisitor returns the rate limiter for the given IP address
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// Allow 10 requests per minute with a burst of 5
		limiter := rate.NewLimiter(rate.Every(time.Minute/10), 5)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time
	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimit is the middleware function to enforce rate limiting
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			logrus.Warn("Rate limit exceeded for IP: ", ip)
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}