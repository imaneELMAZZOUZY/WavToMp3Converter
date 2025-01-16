package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/converter"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
)

var sharedMap = &models.SharedMap{
	Map: make(map[string]models.ConversionConfig),
	Mux: &sync.Mutex{},
}

func main() {

	port:= flag.String("port", ":5000", "HTTP network address")
	flag.Parse()

	dbChan := make(chan models.ConversionRecord, 10)

	go converter.Watch(sharedMap)

	go converter.Process(sharedMap, dbChan)

	go converter.ManageDb(dbChan)

	log.Printf("Starting server on %s", *port)
	err := http.ListenAndServe(*port, routes())
	log.Fatal(err)

}
