package main

import (
	"collector/api"
	"collector/config"
	"context"
	"fmt"
	"net/http"
)

func addConfigMiddleware(cfg config.Config, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v+\n", r)

		ctx := context.WithValue(r.Context(), config.ContextConfig, cfg)
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}

func SetupRoutes(cfg config.Config, s *http.ServeMux) {
	// Handles home path
	s.HandleFunc("/", homeHandler)

	// Handles data files API functions
	s.HandleFunc("/uploads/metadata", addConfigMiddleware(cfg, filesHandler))

	// Handles data files API functions
	s.HandleFunc(fmt.Sprintf("/uploads/content/{%s}", config.ContextHash),
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
