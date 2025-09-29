package handler

import (
	mw "collector/middleware"
	"fmt"
	"i4-lib/config"
	"net/http"
)

func SetupRoutes(cfg config.BaseConfig, mux *http.ServeMux) {
	// Home path handler
	mux.HandleFunc("/", homeHandler)

	// Client handler
	mux.Handle("/clients",
		mw.AdminOrFirstAuthentication(ClientsHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}",
		config.ContextClientCode),
		mw.AdminAuthentication(ClientByIdHandler(), cfg))

	// Machine handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines",
		config.ContextClientCode),
		mw.AdminAuthentication(MachinesHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}",
		config.ContextClientCode,
		config.ContextMachineCode),
		mw.AdminAuthentication(MachineByIdHandler(), cfg))

	// Upload metadata handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/uploads/metadata",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(MetadataHandler(), cfg))

	// Upload content handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/uploads/content/{%s}",
		config.ContextClientCode,
		config.ContextHash),
		mw.AdminOrClientAuthentication(ContentByHashHandler(), cfg))

	// Data handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/data",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(DataHandler(), cfg))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from collector API")
}
