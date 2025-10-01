package middleware

import (
	"i4-lib/service"
	"net/http"
)

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the request
			service.Log("%s request at %s",
				r.Method,
				r.URL,
			)
			h.ServeHTTP(w, r)
		})
	}
}
