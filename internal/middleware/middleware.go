package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// LocalhostOnly restricts access to localhost only
func LocalhostOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.RemoteAddr
		if idx := strings.LastIndex(host, ":"); idx != -1 {
			host = host[:idx]
		}

		fmt.Println()
		
		// Allow localhost, 127.0.0.1, and ::1 (IPv6 localhost)
		if host != "127.0.0.1" && host != "localhost" && host != "::1" && host != "[::1]" {
			http.Error(w, "Access denied: Settings page only accessible from localhost", http.StatusForbidden)
			return
		}
		
		next(w, r)
	}
}

func APIKeyAuth(apiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			providedKey := r.Header.Get("X-API-Key")
			if providedKey == "" {
				http.Error(w, `{"error": "missing API key"}`, http.StatusUnauthorized)
				return
			}
			
			if providedKey != apiKey {
				http.Error(w, `{"error": "invalid API key"}`, http.StatusUnauthorized)
				return
			}
			
			next(w, r)
		}
	}
}