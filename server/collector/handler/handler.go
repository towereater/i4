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

	// Gauge data handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/gauge",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(GaugeHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/gauge/sum",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(GaugeSumHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/gauge/count",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(GaugeCountHandler(), cfg))

	// Interval data handler
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/interval",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(IntervalHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/interval/sum",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(IntervalSumHandler(), cfg))
	mux.Handle(fmt.Sprintf("/clients/{%s}/data/interval/count",
		config.ContextClientCode),
		mw.AdminOrClientAuthentication(IntervalCountHandler(), cfg))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response output
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from collector API")
}
