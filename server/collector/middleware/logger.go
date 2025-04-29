package middleware

import (
	"fmt"
	"net/http"
)

func Logger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		fmt.Printf("Request received:\nMethod: %s\nHeaders: %+v\n",
			r.Method,
			r.Header,
		)

		h.ServeHTTP(w, r)
	})
}
