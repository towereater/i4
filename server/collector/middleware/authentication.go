package middleware

import (
	"context"
	"i4-lib/config"
	"i4-lib/db"
	"i4-lib/service"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func AuthenticateAdminOrClient() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Logging
			service.Log("Headers: %+v", r.Header)

			// Get the api key from request
			apiKey := r.Header.Get("Authentication")
			if apiKey == "" {
				service.Log("Authentication token invalid: %s", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Extract config
			cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err == mongo.ErrNoDocuments {
				service.Log("No client with api key %s", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if err != nil {
				service.Log("Error while searching client with api key %s: %s",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Extract extra parameters
			code := r.PathValue(string(config.ContextClientCode))
			if client.Code != code && client.Code != config.ClientAdminCode {
				service.Log("Client code %s not enabled to manage client %s",
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
			// Get the api key from request
			apiKey := r.Header.Get("Authentication")
			if apiKey == "" {
				service.Log("Authentication token invalid: %s", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Extract config
			cfg := r.Context().Value(config.ContextConfig).(config.BaseConfig)

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err == mongo.ErrNoDocuments {
				service.Log("No client with api key %s", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if err != nil {
				service.Log("Error while searching client with api %s: %s",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if client.Code != config.ClientAdminCode {
				service.Log("Client with api key %s is not admin", apiKey)
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
				service.Log("Error while checking any client existence: %s", err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			// If no client is defined then procede to serve request
			if err == mongo.ErrNoDocuments {
				h.ServeHTTP(w, r)
				return
			}

			// Get the api key from request
			apiKey := r.Header.Get("Authentication")
			if apiKey == "" {
				service.Log("Authentication token invalid: %s", apiKey)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Find the client associated to the given api key
			client, err := db.SelectClientByApiKey(cfg.DB, apiKey)
			if err != nil && err != mongo.ErrNoDocuments {
				service.Log("Error while searching client with api %s: %s",
					apiKey, err.Error())
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if client.Code != config.ClientAdminCode {
				service.Log("Client with api key %s is not admin", apiKey)
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
