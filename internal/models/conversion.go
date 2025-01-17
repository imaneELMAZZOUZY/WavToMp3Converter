package models

import (
	"sync"
)


type SharedMap struct {
	Map map[string]ConversionConfig
	Mux *sync.Mutex
}

// Configuration struct to hold the values from the JSON file
type ConversionConfig struct {
	InputFile  string `json:"input_file"`
	OutputFile string `json:"output_file"`
	Codec      string `json:"codec"`
	Bitrate    string `json:"bitrate"`
	SampleRate string `json:"sample_rate"`
	Channels   string `json:"channels"`
}

// Record struct to hold the values to be inserted into the database
type ConversionRecord struct {
	Id               int
	InputFile        string
	OutputFile       string
	Codec            string
	Bitrate          string
	SampleRate       string
	Channels         string
	ConversionStatus string
	StartTime        string
	EndTime          string
}

type CurrentConfig struct {
	Config ConversionConfig
	StartTime string
}


