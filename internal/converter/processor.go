package converter

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)

var directoryToWatch = flag.String("d", "samples/Watched_folder", "Directory to watch for changes")

var ErrFFMPEGNotFound = errors.New("ffmpeg.exe not found")

var CurrentJobs = &struct {
	Map map[string]models.CurrentConfig
	Mux *sync.Mutex
}{
	Map: make(map[string]models.CurrentConfig),
	Mux: &sync.Mutex{},
}

func Process(j JobConverter, sm *models.SharedMap, dbChan chan<- models.ConversionRecord) {

	// Semaphore to limit the number of concurrent jobs
	semaphore := make(chan struct{}, 5) // max 5

	for {
		if len(sm.Map) > 0 {

			sm.Mux.Lock()

			// Process one item from the map
			for key, value := range sm.Map {

				delete(sm.Map, key)

				semaphore <- struct{}{}

				go func() {
					defer func() {

						time.Sleep(time.Second * 10)

						CurrentJobs.Mux.Lock()
						delete(CurrentJobs.Map, value.InputFile)
						CurrentJobs.Mux.Unlock()

						// Release a slot in the semaphore after job is done
						<-semaphore
					}()

					record, err := j.Run(value)
					if err != nil {
						fmt.Println("error while running conversion process: ", err)
						return
					}
					dbChan <- record
				}()

				break
			}

			sm.Mux.Unlock()

		}

		time.Sleep(time.Second * 2)
	}

}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func runCmd(name string, args []string) error {
	cmd := exec.Command(
		name,
		args...,
	)

	// Run the command
	return cmd.Run()
}

type FileChecker func(path string) bool
type CmdRunner interface {
	Run(name string, args []string) error
}

type JobConverter struct {
	IsFileExist FileChecker
	CmdRunner   CmdRunner
}

func (j JobConverter) Run(jsonConfig models.ConversionConfig) (models.ConversionRecord, error) {
	startTime := time.Now()

	CurrentJobs.Mux.Lock()

	CurrentJobs.Map[jsonConfig.InputFile] = models.CurrentConfig{
		Config:    jsonConfig,
		StartTime: startTime.Format(time.RFC3339),
	}
	CurrentJobs.Mux.Unlock()

	// Define the path to ffmpeg.exe in the assets folder
	ffmpegPath := filepath.Join("bin", "ffmpeg.exe")

	// Check if ffmpeg.exe exists
	if !j.IsFileExist(ffmpegPath) {
		return models.ConversionRecord{}, fmt.Errorf("%w in %s", ErrFFMPEGNotFound, ffmpegPath)
	}

	// Prepare and run the FFmpeg command
	err := j.CmdRunner.Run(
		ffmpegPath,
		[]string{
			"-i", *directoryToWatch + "/" + jsonConfig.InputFile,
			"-codec:a", jsonConfig.Codec,
			"-b:a", jsonConfig.Bitrate,
			"-ar", jsonConfig.SampleRate,
			"-ac", jsonConfig.Channels,
			*directoryToWatch + "/" + jsonConfig.OutputFile,
		},
	)

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

	// Return conversion result
	return conversionRecord, nil
}
