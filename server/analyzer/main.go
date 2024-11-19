package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"analyzer/config"

	"github.com/segmentio/kafka-go"
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
	ctx := context.WithValue(context.Background(), config.ContextConfig, cfg)

	// Main loop
	for {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Queue.Timeout)*time.Second)

		// Poll
		unqueueFile(ctx)

		cancel()
	}
}

func unqueueFile(ctx context.Context) error {
	// Extracting config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Creating topic reader with timeout
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Queue.Brokers,
		GroupID:  cfg.Queue.Uploads.Group,
		Topic:    cfg.Queue.Uploads.Topic,
		MaxBytes: 10e6,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return r.Close()
}

/*

Poll costante della coda in attesa di nuovi dati o in alternativa attesa in linea di dati
Se dei dati sono pronti -> vengono scaricati, convertiti e elaborati
Per ciascuna riga del file, si converte nel formato generico e si analizza il tipo di dato
In base al tipo di dato si esegue un inserimento sulla tabella corretta
Terminata l'elaborazione di un file si procede con quella successiva
Non sono necessarie altre operazioni

*/
