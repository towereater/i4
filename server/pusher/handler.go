package main

import (
	"context"
	"fmt"
	"net/http"
	"pusher/api"
	"pusher/config"
)

func addConfigMiddleware(cfg config.Config, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), config.ContextConfig, cfg)
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}

func SetupRoutes(cfg config.Config, s *http.ServeMux) {
	// Handles home path
	s.HandleFunc("/", homeHandler)

	// Handles data files API functions
	s.HandleFunc("/data-files", addConfigMiddleware(cfg, filesHandler))

	// Handles data files API functions
	s.HandleFunc(fmt.Sprintf("/data-files/{%s}", config.ContextHash),
		addConfigMiddleware(cfg, filesByIdHandler))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from files API")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertMetadata(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func filesByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertFile(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}
