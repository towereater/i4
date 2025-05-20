package middleware

import (
	"fmt"
	"net/http"
)

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the request
			fmt.Printf("Request received:\nMethod: %s\nHeaders: %+v\n",
				r.Method,
				r.Header,
			)
			h.ServeHTTP(w, r)
		})
	}
}
