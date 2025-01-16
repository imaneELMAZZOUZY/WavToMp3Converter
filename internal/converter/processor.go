package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"sync"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/joho/godotenv"
)

var CurrentJobs = &models.SharedMap{
	Map: make(map[string]models.ConversionConfig),
	Mux: &sync.Mutex{},
}

func Process(sm *models.SharedMap, dbChan chan<- models.ConversionRecord) {

	// Semaphore to limit the number of concurrent jobs
	semaphore := make(chan struct{}, 5) // max 5

	for {
		if len(sm.Map) > 0 {

			sm.Mux.Lock()

			// Process one item from the map
			for key, value := range sm.Map {

				delete(sm.Map, key)

				CurrentJobs.Mux.Lock()
				CurrentJobs.Map[value.InputFile] = value
				CurrentJobs.Mux.Unlock()

				// Acquire a semaphore before starting a job
				semaphore <- struct{}{}

				go jobConverter(value, dbChan, semaphore)

				break
			}

			sm.Mux.Unlock()

		}

		time.Sleep(time.Second * 10) 
	}

}

func jobConverter(jsonConfig models.ConversionConfig, dbChan chan<- models.ConversionRecord, semaphore chan struct{}) {

	defer func() {

		time.Sleep(time.Second * 10)

		CurrentJobs.Mux.Lock()
		delete(CurrentJobs.Map, jsonConfig.InputFile)
		CurrentJobs.Mux.Unlock()

		// Release a slot in the semaphore after job is done
		<-semaphore
	}()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	directoryToWatch := os.Getenv("DIRECTORY_TO_WATCH")

	// Define the path to ffmpeg.exe in the assets folder
	ffmpegPath := filepath.Join("bin", "ffmpeg.exe")

	// Check if ffmpeg.exe exists
	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
		fmt.Printf("ffmpeg.exe not found in %s\n", ffmpegPath)
		return
	}

	startTime := time.Now()
	// Prepare the FFmpeg command
	cmd := exec.Command(
		ffmpegPath,
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
		conversionStatus = "failed"
	} else {
		conversionStatus = "successful"
		fmt.Println("Conversion successful! Output file:", jsonConfig.OutputFile)
	}

	// Prepare conversion record
	conversionRecord := models.ConversionRecord{
		InputFile:        jsonConfig.InputFile,
		OutputFile:       jsonConfig.OutputFile,
		Codec:            jsonConfig.Codec,
		Bitrate:          jsonConfig.Bitrate,
		SampleRate:       jsonConfig.SampleRate,
		Channels:         jsonConfig.Channels,
		ConversionStatus: conversionStatus,
		StartTime:        startTime.Format(time.RFC3339),
		EndTime:          time.Now().Format(time.RFC3339),
	}

	// Send to channel
	dbChan <- conversionRecord

}
