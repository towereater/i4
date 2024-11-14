package main

import (
	"fmt"
	"os"

	"worker/config"
)

func main() {
	// Get run args
	if len((os.Args)) < 2 {
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

	// Starting up
	fmt.Println(cfg)
}

/*

Poll costante della coda in attesa di nuovi dati o in alternativa attesa in linea di dati
Se dei dati sono pronti -> vengono scaricati, convertiti e elaborati
Per ciascuna riga del file, si converte nel formato generico e si analizza il tipo di dato
In base al tipo di dato si esegue un inserimento sulla tabella corretta
Terminata l'elaborazione di un file si procede con quella successiva
Non sono necessarie altre operazioni

*/
