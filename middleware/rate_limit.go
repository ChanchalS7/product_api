package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func init() {
	go cleanupVisitors()
}

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

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		v, exists := visitors[ip]

		if !exists {
			limiter := rate.NewLimiter(1, 3)
			visitors[ip] = &visitor{limiter, time.Now()}
		} else {
			v.lastSeen = time.Now()
		}
		if !v.limiter.Allow() {
			mu.Unlock()
			logrus.Warn("Rate limit exceeded for this IP:", ip)
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
