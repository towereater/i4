package middleware

import (
	"context"
	"fmt"
	"i4-lib/config"
	"i4-lib/db"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func AuthenticateAdminOrClient() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the authentication header from request
			auth := r.Header["Authentication"]
			if len(auth) != 1 {
				fmt.Printf("Authorization token invalid: %s\n", auth)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Get the api key
			apiKey := auth[0]

			// Extract config
			cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err == mongo.ErrNoDocuments {
				fmt.Printf("No client with api key %s\n", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if err != nil {
				fmt.Printf("Error while searching client with api key %s: %s\n",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Extract extra parameters
			code := r.PathValue(string(config.ContextClientCode))
			if client.Code != code && client.Code != config.ClientAdminCode {
				fmt.Printf("Client code %s not enabled to manage client %s\n",
					client.Code, code)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Add user code to the request context
			ctx := context.WithValue(r.Context(), config.ContextClientCode, client.Code)
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}

func AuthenticateAdmin() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the authentication header from request
			auth := r.Header["Authentication"]
			if len(auth) != 1 {
				fmt.Printf("Auth token invalid: %s\n", auth)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Get the api key
			apiKey := auth[0]

			// Extract config
			cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err == mongo.ErrNoDocuments {
				fmt.Printf("No client with api key %s\n", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if err != nil {
				fmt.Printf("Error while searching client with api %s: %s\n",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if client.Code != config.ClientAdminCode {
				fmt.Printf("Client with api key %s is not admin\n", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Add user code to the request context
			ctx := context.WithValue(r.Context(), config.ContextClientCode, client.Code)
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}

func AuthenticateAdminOrFirstUser() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract config
			cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

			// Check if no client is defined at all
			_, err := db.SelectAnyClient(cfg.DB)
			if err != nil && err != mongo.ErrNoDocuments {
				fmt.Printf("Error while checking any client existence: %s\n", err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			// If no client is defined then procede to serve request
			if err == mongo.ErrNoDocuments {
				h.ServeHTTP(w, r)
				return
			}

			// Get the authentication header from request
			auth := r.Header["Authentication"]
			if len(auth) != 1 {
				fmt.Printf("Auth token invalid: %s\n", auth)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Get the api key
			apiKey := auth[0]

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err != nil && err != mongo.ErrNoDocuments {
				fmt.Printf("Error while searching client with api %s: %s\n",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if client.Code != config.ClientAdminCode {
				fmt.Printf("Client with api key %s is not admin\n", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Add user code to the request context
			ctx := context.WithValue(r.Context(), config.ContextClientCode, client.Code)
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}
