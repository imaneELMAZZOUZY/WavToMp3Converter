package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/joho/godotenv"
)

// monitor a directory for .json and .wav files and update the shared map.
func Watch(sm *models.SharedMap) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	directoryToWatch := os.Getenv("DIRECTORY_TO_WATCH")
	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error while creating a watcher : ", err)
	}

	defer watcher.Close()

	err = watcher.Add(directoryToWatch)
	if err != nil {
		fmt.Println("Error adding directory to watcher:", err)
	}

	fmt.Printf("Watching directory: %s\n", directoryToWatch)

	// Track file events by base filename
	fileCreationCount := make(map[string]int)
	for {
		select {
		case event, ok := <-watcher.Events:
			// if the channel is empty and closed
			if !ok {
				return
			}

			// Handle creation events
			if event.Op&fsnotify.Create == fsnotify.Create {

				// getting the file name and extension
				filebase := filepath.Base(event.Name)
				filename, ext := strings.Split(filebase, ".")[0], strings.Split(filebase, ".")[1]

				// Track the number of files created with the same base name (ensure .wav and .json are present)
				if ext == "json" || ext == "wav" {
					fileCreationCount[filename]++
				}

				// When both .wav and .json files are created, process them
				if fileCreationCount[filename] == 2 {
					config, err := jsonToStruct(directoryToWatch + "/" + filename + ".json")
					if err != nil {
						fmt.Println(err)
					}

					// Update shared map with conversion configuration
					sm.Mux.Lock()
					sm.Map[filename] = config
					sm.Mux.Unlock()

					delete(fileCreationCount, filename)
				}

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
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
