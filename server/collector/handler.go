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
		ctx := context.WithValue(r.Context(), config.ContextConfig, cfg)
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}

func setupRoutes(cfg config.Config, s *http.ServeMux) {
	// Home path handler
	s.HandleFunc("/", homeHandler)

	// Upload metadata handler
	s.HandleFunc("/uploads/metadata", addConfigMiddleware(cfg, metadataHandler))

	// Upload content handler
	s.HandleFunc(fmt.Sprintf("/uploads/content/{%s}", config.ContextHash),
		addConfigMiddleware(cfg, contentByIdHandler))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from files API")
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertMetadata(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func contentByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Check of the method request
	switch r.Method {
	case "POST":
		api.InsertContent(w, r)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}
