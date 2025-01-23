package main

import ( 
	"path/filepath"
	"strings"
	"github.com/fsnotify/fsnotify"
)

// Monitor a directory for .json and .wav files and update the shared map.
func (appDep *appDep) Watch() {
	
	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		appDep.logger.Error("Error while creating a watcher : ", "error", err)
	}

	defer watcher.Close()

	err = watcher.Add(*DirectoryToWatch)
	if err != nil {
		appDep.logger.Error("Error adding directory to watcher:", "error", err)
		return
	}

	appDep.logger.Info("Watching...", "directory", *DirectoryToWatch)

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
					appDep.logger.Warn("File does not have an extension:", "filename", filebase)
					continue
				 }
				 filename := strings.TrimSuffix(filebase, ext)

				// Track the number of files created with the same base name (ensure .wav and .json are present)
				if ext == ".json" || ext == ".wav" {
					fileCreationCount[filename]++
				}

				// When both .wav and .json files are created, process them
				if fileCreationCount[filename] == 2 {
					config, err := jsonToStruct(*DirectoryToWatch + "/" + filename + ".json")
					if err != nil {
						appDep.logger.Error(err.Error())
					}

					// Update shared map with conversion configuration
					appDep.sharedMap.Mux.Lock()
					appDep.sharedMap.Map[filename] = config
					appDep.sharedMap.Mux.Unlock()

					delete(fileCreationCount, filename)
				}	
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			appDep.logger.Error("Watcher", "error", err)
		}
	}
}


