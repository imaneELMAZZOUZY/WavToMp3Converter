package converter

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/joho/godotenv"
)

func Process(sm *models.SharedMap, dbChan chan<- models.ConversionRecord) {

	// Semaphore to limit the number of concurrent jobs
	semaphore := make(chan struct{}, 5) // max 5

	for {
		if len(sm.Map) > 0 {

			sm.Mux.Lock()

			// process one item from the map
			for key, value := range sm.Map {

				delete(sm.Map, key)

				// Acquire a semaphore before starting a job
				semaphore <- struct{}{}

				go jobConverter(value, dbChan, semaphore)

				break
			}
			sm.Mux.Unlock()

		}

		time.Sleep(time.Millisecond * 20)
	}

}

func jobConverter(jsonConfig models.ConversionConfig, dbChan chan<- models.ConversionRecord, semaphore chan struct{}) {

	defer func() {
		// Release a slot in the semaphore after job is done
		<-semaphore
	}()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	directoryToWatch := os.Getenv("DIRECTORY_TO_WATCH")

	startTime := time.Now()
	// Prepare the FFmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-i", directoryToWatch+"/"+jsonConfig.InputFile,
		"-codec:a", jsonConfig.Codec,
		"-b:a", jsonConfig.Bitrate,
		"-ar", jsonConfig.SampleRate,
		"-ac", jsonConfig.Channels,
		directoryToWatch+"/"+jsonConfig.OutputFile,
	)

	// Run the command
	err = cmd.Run()

	var conversionStatus string
	if err != nil {
		fmt.Println("Error executing FFmpeg command:", err)
		conversionStatus = "Failed"
	} else {
		conversionStatus = "Successful"
		fmt.Println("Conversion successful! Output file:", jsonConfig.OutputFile)
	}

	duration := time.Since(startTime).String()

	// Prepare conversion record
	conversionRecord := models.ConversionRecord{
		InputFile:        jsonConfig.InputFile,
		OutputFile:       jsonConfig.OutputFile,
		Codec:            jsonConfig.Codec,
		Bitrate:          jsonConfig.Bitrate,
		SampleRate:       jsonConfig.SampleRate,
		Channels:         jsonConfig.Channels,
		ConversionStatus: conversionStatus,
		Duration:         duration,
		Timestamp:        time.Now().Format(time.RFC3339),
	}

	// send to channel
	dbChan <- conversionRecord

}
