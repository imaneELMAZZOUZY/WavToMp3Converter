package main

import (
	"net/http"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/converter"
	"strings"
	"encoding/json"
)


func currentJobsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(converter.CurrentJobs.Map)
	if err != nil {
		http.Error(w, "Error encoding current jobs to JSON", http.StatusInternalServerError)
	}

	
}

func waitingJobsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(converter.CurrentJobs.Map)
	if err != nil {
		http.Error(w, "Error encoding waiting jobs to JSON", http.StatusInternalServerError)
	}

}

func finishedJobsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	status := r.URL.Query().Get("status")
	status = strings.ToLower(status)
	if status != "successful" && status != "failed" && status != "" {

		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	finishedJobs, err:= converter.GetRecords(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the finished jobs to JSON and send it to the response writer
	err = json.NewEncoder(w).Encode(finishedJobs)
	if err != nil {
		http.Error(w, "Error encoding finished jobs to JSON", http.StatusInternalServerError)
	}


}