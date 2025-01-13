package main

import (
	"sync"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/converter"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)

func main() {

	done := make(chan bool)

	dbChan := make(chan models.ConversionRecord, 10)

	sharedMap := &models.SharedMap{
		Map: make(map[string]models.ConversionConfig),
		Mux: &sync.Mutex{},
	}

	go converter.Watch(sharedMap)

	go converter.Process(sharedMap, dbChan)

	go converter.ManageDb(dbChan)

	<-done

}
