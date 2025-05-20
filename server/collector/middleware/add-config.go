package middleware

import (
	"collector/config"
	"context"
	"net/http"
)

func AddConfig(cfg config.Config) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add a basic configuration to the request context
			ctx := context.WithValue(r.Context(), config.ContextConfig, cfg)
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}
