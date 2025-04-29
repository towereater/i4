package handler

import (
	"collector/config"
	"collector/middleware"
	"fmt"
	"net/http"
)

func SetupRoutes(cfg config.Config, mux *http.ServeMux) {
	// Home path handler
	mux.HandleFunc("/", homeHandler)

	// Client handler
	mux.HandleFunc("/clients",
		middleware.Logger(middleware.AddConfig(cfg, middleware.AuthenticateAdmin(ClientsHandler))))
	mux.HandleFunc(fmt.Sprintf("/clients/{%s}", config.ContextClientCode),
		middleware.Logger(middleware.AddConfig(cfg, middleware.AuthenticateClient(ClientsByIdHandler))))

	// Machine handler
	mux.HandleFunc(fmt.Sprintf("/clients/{%s}/machines", config.ContextClientCode),
		middleware.Logger(middleware.AddConfig(cfg, middleware.AuthenticateAdmin(MachinesHandler))))
	mux.HandleFunc(fmt.Sprintf("/clients/{%s}/machines/{%s}", config.ContextClientCode, config.ContextMachineCode),
		middleware.Logger(middleware.AddConfig(cfg, middleware.AuthenticateAdmin(MachinesByIdHandler))))

	// Upload metadata handler
	mux.HandleFunc("/uploads/metadata",
		middleware.AddConfig(cfg, middleware.AuthenticateClient(MetadataHandler)))

	// Upload content handler
	mux.HandleFunc(fmt.Sprintf("/uploads/content/{%s}", config.ContextHash),
		middleware.AddConfig(cfg, middleware.AuthenticateClient(ContentHandler)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from files API")
}
