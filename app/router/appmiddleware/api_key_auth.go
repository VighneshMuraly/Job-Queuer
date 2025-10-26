package appmiddleware

import (
	"fmt"
	"job-queuer/app/util/env"
	"net/http"

	"github.com/gorilla/mux"
)

func APIKeyAuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("Authorization")
			if apiKey == "" {
				// also allow ?api_key=... as fallback
				apiKey = r.URL.Query().Get("api_key")
			}
			validAPIKey := env.GetAPIKey()
			if apiKey != validAPIKey {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Unauthorized: Invalid or missing API key\n")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
