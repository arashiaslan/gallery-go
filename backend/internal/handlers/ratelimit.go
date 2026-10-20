package handlers

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	authWindow        = 15 * time.Minute
	authMaxAttempts   = 10
	authLimiterMaxKey = 4096
)

type authRateWindow struct {
	start time.Time
	count int
}

type authRateLimiter struct {
	mu      sync.Mutex
	windows map[string]*authRateWindow
}

func newAuthRateLimiter() *authRateLimiter {
	return &authRateLimiter{windows: make(map[string]*authRateWindow)}
}

// Allow reports whether the key may proceed, consuming one attempt.
func (l *authRateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if len(l.windows) >= authLimiterMaxKey {
		for k, w := range l.windows {
			if now.Sub(w.start) > authWindow {
				delete(l.windows, k)
			}
		}
	}

	w, ok := l.windows[key]
	if !ok || now.Sub(w.start) > authWindow {
		l.windows[key] = &authRateWindow{start: now, count: 1}
		return true
	}
	w.count++
	return w.count <= authMaxAttempts
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
