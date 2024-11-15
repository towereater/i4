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

	// Main loop
	for {
		// Poll
	}
}

/*
func unqueueFile(ctx context.Context) error {
	// Extracting config
	cfg := ctx.Value(config.ContextConfig).(config.Config)

	// Creating topic reader with timeout
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		GroupID:  "consumer-group-id",
		Topic:    "topic-A",
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	/////////
	w := &kafka.Writer{
		Addr:  kafka.TCP(fmt.Sprintf("%s:%s", cfg.Queue.Host, cfg.Queue.Port)),
		Topic: cfg.Queue.Topic,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Queue.Timeout)*time.Second)
	defer cancel()

	// Preparing data for queue
	datetime := time.Now().Format(time.DateTime)
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, hash)

	// Writing hash on queue
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(datetime),
			Value: h,
		},
	)
	if err != nil {
		return err
	}

	return w.Close()
}

*/

/*

Poll costante della coda in attesa di nuovi dati o in alternativa attesa in linea di dati
Se dei dati sono pronti -> vengono scaricati, convertiti e elaborati
Per ciascuna riga del file, si converte nel formato generico e si analizza il tipo di dato
In base al tipo di dato si esegue un inserimento sulla tabella corretta
Terminata l'elaborazione di un file si procede con quella successiva
Non sono necessarie altre operazioni

*/
