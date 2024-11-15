package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"collector/config"
)

func main() {
	// Get run args
	if len(os.Args) < 2 {
		println("No config file set")
		os.Exit(1)
	}
	configPath := os.Args[1]

	// Setup machine config
	fmt.Println("Loading configuration")
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		println("Error while reading config file:", err)
		os.Exit(2)
	}

	// Creation of the mux
	mux := http.NewServeMux()

	// Setup server routes
	fmt.Println("Setting up routes")
	SetupRoutes(cfg, mux)

	// Creation of the server
	server := &http.Server{
		Handler: mux,
	}
	ln, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		println("Error while assigning server port:", err)
		os.Exit(3)
	}

	// Starting up
	fmt.Println("Ready to listen incoming requests")
	server.Serve(ln)
	if err != nil {
		println("Error while starting up server:", err)
		os.Exit(4)
	}
}
