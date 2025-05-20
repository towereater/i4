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
		mw.LoggedAdminAuthentication(ClientsHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}",
		config.ContextClientCode),
		mw.LoggedAdminAuthentication(ClientsByIdHandler(), cfg))

	// Machine handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines",
		config.ContextClientCode),
		mw.LoggedAdminAuthentication(MachinesHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}",
		config.ContextClientCode,
		config.ContextMachineCode),
		mw.LoggedAdminAuthentication(MachinesByIdHandler(), cfg))

	// Upload metadata handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}/uploads/metadata",
		config.ContextClientCode,
		config.ContextMachineCode),
		mw.LoggedClientAuthentication(MetadataHandler(), cfg))

	// Upload content handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/machines/{%s}/uploads/content/{%s}",
		config.ContextClientCode,
		config.ContextMachineCode,
		config.ContextHash),
		mw.LoggedClientAuthentication(ContentHandler(), cfg))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from files API")
}
