package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func jsonToStruct(filepath string) (models.ConversionConfig, error) {
	var config models.ConversionConfig

	// Attempt to open the file with retries
	for retries := 0; retries < 10; retries++ {
		file, err := os.Open(filepath)
		if err == nil {
			defer file.Close()

			// Decode the JSON content into the Config struct
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&config); err != nil {
				return config, fmt.Errorf("error decoding JSON: %w", err)
			}
			return config, nil
		}

		// Wait a bit before retrying
		time.Sleep(time.Millisecond * 20)
	}

	return config, fmt.Errorf("failed to open file after many retries")
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wavtomp3converter_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}
