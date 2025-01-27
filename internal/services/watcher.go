package services

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/app"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/helpers"
)

// Monitor a directory for .json and .wav files and update the shared map.
func  Watch(appDep *app.AppDep) {
	
	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		appDep.Logger.Error("Error while creating a watcher : ", "error", err)
	}

	defer watcher.Close()

	err = watcher.Add(*app.DirectoryToWatch)
	if err != nil {
		appDep.Logger.Error("Error adding directory to watcher:", "error", err)
		return
	}

	appDep.Logger.Info("Watching...", "directory", *app.DirectoryToWatch)

	// Track file events by base filename
	fileCreationCount := make(map[string]int)
	for {
		select {
		case event, ok := <-watcher.Events:
			// If the channel is empty and closed
			if !ok {
				return
			}

			// Handle creation events
			if event.Op&fsnotify.Create == fsnotify.Create {

				// Getting the file name and extension
				filebase := filepath.Base(event.Name)
				ext := filepath.Ext(filebase)
				 if ext == "" {
					appDep.Logger.Warn("File does not have an extension:", "filename", filebase)
					continue
				 }
				 filename := strings.TrimSuffix(filebase, ext)

				// Track the number of files created with the same base name (ensure .wav and .json are present)
				if ext == ".json" || ext == ".wav" {
					fileCreationCount[filename]++
				}

				// When both .wav and .json files are created, process them
				if fileCreationCount[filename] == 2 {
					config, err := helpers.JsonToStruct(*app.DirectoryToWatch + "/" + filename + ".json")
					if err != nil {
						appDep.Logger.Error(err.Error())
					}

					// Update shared map with conversion configuration
					
					appDep.SharedMap.Store(filename, config) 
					

					delete(fileCreationCount, filename)
				}	
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			appDep.Logger.Error("Watcher", "error", err)
		}
	}
}


