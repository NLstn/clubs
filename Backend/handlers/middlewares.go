package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/NLstn/clubs/auth"
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
		// SECURITY: Use FRONTEND_URL from environment instead of wildcard (*)
		// Wildcard (*) with credentials is a security vulnerability
		allowedOrigin := os.Getenv("FRONTEND_URL")
		if allowedOrigin == "" {
			// Fallback to localhost for development if not set
			allowedOrigin = "http://localhost:5173"
		}
		
		// Check if origin is allowed
		origin := r.Header.Get("Origin")
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Refresh-Token, X-API-Key")

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
			// SECURITY: Extract IP address carefully
			// In production with a reverse proxy (e.g., nginx, load balancer),
			// X-Forwarded-For should be used. In development, use RemoteAddr.
			// We use X-Real-IP as it's typically set by trusted reverse proxies
			// and is harder to spoof than X-Forwarded-For
			var ip string
			
			// Try X-Real-IP first (set by trusted reverse proxy)
			ip = r.Header.Get("X-Real-IP")
			if ip == "" {
				// Fallback to X-Forwarded-For, taking only the first (client) IP
				xff := r.Header.Get("X-Forwarded-For")
				if xff != "" {
					// X-Forwarded-For may contain multiple IPs: "client, proxy1, proxy2"
					// Take the first one (client IP)
					ips := strings.Split(xff, ",")
					ip = strings.TrimSpace(ips[0])
				}
			}
			
			// Fallback to RemoteAddr if no proxy headers
			if ip == "" {
				// RemoteAddr format is "IP:port", extract just the IP
				addr := r.RemoteAddr
				if idx := strings.LastIndex(addr, ":"); idx != -1 {
					ip = addr[:idx]
				} else {
					ip = addr
				}
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

// APIKeyAuthMiddleware validates API keys from X-API-Key or Authorization: ApiKey headers
func APIKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for API key in X-API-Key header (preferred)
		apiKey := r.Header.Get("X-API-Key")

		// If not found, check Authorization header for "ApiKey" scheme
		if apiKey == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		// If no API key found, return unauthorized
		if apiKey == "" {
			log.Println("API key authentication failed: no API key provided")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the API key
		userID, _, err := auth.ValidateAPIKey(apiKey)
		if err != nil {
			log.Printf("API key authentication failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set user ID in context
		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CompositeAuthMiddleware tries JWT Bearer token first, then API key authentication
func CompositeAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Try JWT Bearer token authentication first
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Use the existing JWT auth logic
			auth.AuthMiddleware(next).ServeHTTP(w, r)
			return
		}

		// Try API key authentication
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" && strings.HasPrefix(authHeader, "ApiKey ") {
			apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
		}

		if apiKey != "" {
			// Validate the API key
			userID, _, err := auth.ValidateAPIKey(apiKey)
			if err != nil {
				log.Printf("Authentication failed: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set user ID in context
			ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// No valid authentication found
		log.Println("Authentication failed: no valid credentials provided")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
