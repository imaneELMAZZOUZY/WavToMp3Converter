package services

import (
	"fmt"
	"os"

	"os/exec"
	"path/filepath"

	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/app"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/helpers"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)



func Process(appDep *app.AppDep, j JobConverter) {

	// Semaphore to limit the number of concurrent jobs
	semaphore := make(chan struct{}, 5) // max 5

	for {

		appDep.SharedMap.Range(func(k, val any) bool {
			if k == nil {
				return false
			}
			appDep.SharedMap.Delete(k)

			semaphore <- struct{}{}

			value := val.(models.ConversionConfig)
			key := k.(string)

			go func() {

				defer func() {

					time.Sleep(time.Second * 10)

					appDep.CurrentJobs.Delete(value.InputFile)

					// Release a slot in the semaphore after job is done
					<-semaphore
				}()

				startTime := time.Now().Format(time.RFC3339)

				appDep.CurrentJobs.Store(value.InputFile, models.CurrentConfig{
					Config:    value,
					StartTime: startTime,
				})

				record, err := j.Run(value, startTime)
				if err != nil {
					appDep.Logger.Error("error while running conversion process :", "error", err)
				} else {
					appDep.Logger.Info("conversion successful!", "input_file", record.InputFile, "output_file", record.OutputFile)
					appDep.DbChan <- record

					err := os.Remove(*app.DirectoryToWatch + "/" + key + ".json")  
					if err != nil {
						appDep.Logger.Error("error while removing json file :", "error", err)
					}

					err = os.Remove(*app.DirectoryToWatch + "/" + key + ".wav")  
					if err != nil {
						appDep.Logger.Error("error while removing wav file :", "error", err)
					}
				}

			}()

			return false

		})

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
		IsFileExist: helpers.IsFileExist,
		CmdRunner:   cmdRunner{},
	}
}

func (j JobConverter) Run(jsonConfig models.ConversionConfig, startTime string) (models.ConversionRecord, error) {

	// Define the path to ffmpeg.exe in the assets folder
	ffmpegPath := filepath.Join("bin", "ffmpeg.exe")

	// Check if ffmpeg.exe exists
	if !j.IsFileExist(ffmpegPath) {
		return models.ConversionRecord{}, fmt.Errorf("%w in %s", app.ErrFFMPEGNotFound, ffmpegPath)
	}

	// Run the command
	err := j.CmdRunner.Run(ffmpegPath, []string{
		"-i",
		*app.DirectoryToWatch + "/" + jsonConfig.InputFile,
		"-codec:a", jsonConfig.Codec,
		"-b:a", jsonConfig.Bitrate,
		"-ar", jsonConfig.SampleRate,
		"-ac", jsonConfig.Channels,
		*app.OutputDirectory + "/" + jsonConfig.OutputFile,
	})

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

	return conversionRecord, nil

}
