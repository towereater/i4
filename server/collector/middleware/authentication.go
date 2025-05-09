package middleware

import (
	"collector/config"
	"collector/db"
	"context"
	"fmt"
	"net/http"
)

func AuthenticateClient(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authentication header from request
		auth := r.Header["Authentication"]
		if len(auth) != 1 {
			// Write response output
			fmt.Printf("Authorization token invalid: %s\n", auth)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Get the api key
		apiKey := auth[0]

		// Find the client associated to the given api key
		client, err := db.SelectClientByApiKey(r.Context(), apiKey)
		if err != nil {
			// Write response output
			fmt.Printf("Error while searching client with api %s: %s\n", apiKey, err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Extract extra parameters
		code := r.PathValue(string(config.ContextClientCode))
		if client.Code != code && client.Code != "00000" {
			// Write response output
			fmt.Printf("Client code %s not enabled to manage client %s\n", client.Code, code)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Add user code to the request context
		ctx := context.WithValue(r.Context(), config.ContextClientCode, client.Code)
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}

func AuthenticateAdmin(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authentication header from request
		auth := r.Header["Authentication"]
		if len(auth) != 1 {
			// Write response output
			fmt.Printf("Auth token invalid: %s\n", auth)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Get the api key
		apiKey := auth[0]

		// Find the client associated to the given api key
		client, err := db.SelectClientByApiKey(r.Context(), apiKey)
		if err != nil {
			// Write response output
			fmt.Printf("Error while searching client with api %s: %s\n", apiKey, err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if client.Code != "00000" {
			// Write response output
			fmt.Printf("Client with api %s is not admin\n", apiKey)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Add user code to the request context
		ctx := context.WithValue(r.Context(), config.ContextClientCode, client.Code)
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}
