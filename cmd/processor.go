package main

import (
	"fmt"
	
	"os/exec"
	"path/filepath"
	
	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)





func (appDep *appDep) Process(j JobConverter) {

	// Semaphore to limit the number of concurrent jobs
	semaphore := make(chan struct{}, 5) // max 5

	for {
		if len(appDep.sharedMap.Map) > 0 {

			appDep.sharedMap.Mux.Lock()

			// Process one item from the map
			for key, value := range appDep.sharedMap.Map {

				delete(appDep.sharedMap.Map, key)

				semaphore <- struct{}{}

				go func() {

					defer func() {

						time.Sleep(time.Second * 10)

						appDep.CurrentJobs.Mux.Lock()
						delete(appDep.CurrentJobs.Map, value.InputFile)
						appDep.CurrentJobs.Mux.Unlock()

						// Release a slot in the semaphore after job is done
						<-semaphore
					}()

					startTime := time.Now().Format(time.RFC3339)

					appDep.CurrentJobs.Mux.Lock()

					appDep.CurrentJobs.Map[value.InputFile] = models.CurrentConfig{
						Config:    value,
						StartTime: startTime,
					}
					appDep.CurrentJobs.Mux.Unlock()

					record , err := j.Run(value,startTime)
					if err != nil {
						appDep.logger.Error("error while running conversion process :", "error", err)
					} else {
						appDep.logger.Info("conversion successful!", "input_file", record.InputFile, "output_file", record.OutputFile)
						appDep.dbChan <- record
					}
				
					

				}()

				break
			}

			appDep.sharedMap.Mux.Unlock()

		}

		time.Sleep(time.Second * 2)
	}

}

type cmdRunner struct{}

func (cmdRunner) Run(name string, args []string) error {
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

func NewJobConverter() JobConverter {
	return JobConverter{
		IsFileExist: isFileExist,
		CmdRunner:   cmdRunner{},
	}
}

func (j JobConverter) Run(jsonConfig models.ConversionConfig, startTime string) (models.ConversionRecord, error){

	// Define the path to ffmpeg.exe in the assets folder
	ffmpegPath := filepath.Join("bin", "ffmpeg.exe")

	// Check if ffmpeg.exe exists
	if !j.IsFileExist(ffmpegPath) {
		return models.ConversionRecord{}, fmt.Errorf("%w in %s", ErrFFMPEGNotFound, ffmpegPath)
	}

	// Run the command
	err := j.CmdRunner.Run(ffmpegPath,[]string{
		"-i", 
		*DirectoryToWatch+"/"+jsonConfig.InputFile,
		"-codec:a", jsonConfig.Codec,
		"-b:a", jsonConfig.Bitrate,
		"-ar", jsonConfig.SampleRate,
		"-ac", jsonConfig.Channels,
		*DirectoryToWatch+"/"+jsonConfig.OutputFile,
       } )

	var conversionStatus string
	if err != nil {
		conversionStatus = "failed"

	} else {
		conversionStatus = "successful"
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
		StartTime:        startTime,
		EndTime:          time.Now().Format(time.RFC3339),
	}

	return conversionRecord,nil

}
