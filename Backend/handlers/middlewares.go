package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Refresh-Token")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create custom response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // Default status
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"endpoint=%s method=%s status=%d duration=%v",
			r.URL.Path,
			r.Method,
			rw.status,
			duration,
		)
	})
}

// IPRateLimiter stores rate limiters for IP addresses
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	ttl    time.Duration
	lastOp map[string]time.Time
}

// NewIPRateLimiter creates a new rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		ttl:    time.Hour, // Clean up entries after 1 hour of inactivity
		lastOp: make(map[string]time.Time),
	}

	// Start cleanup routine
	go i.cleanup()
	return i
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	now := time.Now()
	i.lastOp[ip] = now

	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

// cleanup removes old entries periodically
func (i *IPRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute) // Check every minute
		i.mu.Lock()
		now := time.Now()
		for ip, lastOp := range i.lastOp {
			if now.Sub(lastOp) > i.ttl {
				delete(i.ips, ip)
				delete(i.lastOp, ip)
			}
		}
		i.mu.Unlock()
	}
}

// Global rate limiters for different endpoints
var (
	// Strict limiter for authentication endpoints (5 requests per minute)
	authLimiter = NewIPRateLimiter(rate.Limit(5.0/60.0), 5)

	// More lenient limiter for general API endpoints (30 requests per 5 seconds)
	apiLimiter = NewIPRateLimiter(rate.Limit(30.0/5.0), 30)
)

// RateLimitMiddleware creates a middleware with the specified limiter
func RateLimitMiddleware(limiter *IPRateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get IP from X-Forwarded-For header, or fallback to RemoteAddr
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}

			limiter := limiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
