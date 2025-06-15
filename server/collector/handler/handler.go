package handler

import (
	"collector/config"
	mw "collector/middleware"
	"fmt"
	"net/http"
)

func SetupRoutes(cfg config.Config, mux *http.ServeMux) {
	// Home path handler
	mux.HandleFunc("/", homeHandler)

	// Client handler
	mux.Handle("/clients",
		mw.AdminOrFirstAuthentication(ClientsHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}",
		config.ContextClientCode),
		mw.AdminAuthentication(ClientsByIdHandler(), cfg))

	// Machine handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines",
		config.ContextClientCode),
		mw.AdminAuthentication(MachinesHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}",
		config.ContextClientCode,
		config.ContextMachineCode),
		mw.AdminAuthentication(MachinesByIdHandler(), cfg))

	// Upload metadata handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}/metadata",
		config.ContextClientCode,
		config.ContextMachineCode),
		mw.AdminOrClientAuthentication(MetadataHandler(), cfg))

	// Upload content handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}/content/{%s}",
		config.ContextClientCode,
		config.ContextMachineCode,
		config.ContextHash),
		mw.AdminOrClientAuthentication(ContentHandler(), cfg))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from collector API")
}
