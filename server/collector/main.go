package main

import (
	"net"
	"net/http"
	"os"

	"collector/handler"
	"i4-lib/config"
	"i4-lib/service"
)

func main() {
	// Get run args
	if len(os.Args) < 2 {
		service.Log("No config file set\n")
		os.Exit(1)
	}
	configPath := os.Args[1]
	service.Log("Loading configuration from %s\n", configPath)

	// Setup machine config
	var cfg config.BaseConfig
	err := config.LoadConfig(configPath, &cfg)
	if err != nil {
		service.Log("Error while reading config file: %s\n", err.Error())
		os.Exit(2)
	}
	service.Log("Configuration loaded: %+v\n", cfg)

	// Create the mux
	mux := http.NewServeMux()

	// Setup server routes
	service.Log("Setting up routes\n")
	handler.SetupRoutes(cfg, mux)

	// Create the server
	server := &http.Server{
		Handler: mux,
	}
	ln, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		service.Log("Error while assigning server port: %s\n", err.Error())
		os.Exit(3)
	}

	// Starting up
	service.Log("Ready to listen incoming requests\n")
	server.Serve(ln)
	if err != nil {
		service.Log("Error while starting up server: %s\n", err.Error())
		os.Exit(4)
	}
}
