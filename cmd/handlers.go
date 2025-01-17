package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/converter"
)

// currentJobsHandler handles the request to get the current jobs.
// It encodes the current jobs to JSON and writes it to the response.
func currentJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(converter.CurrentJobs.Map)
	if err != nil {
		http.Error(w, "Error encoding current jobs to JSON", http.StatusInternalServerError)
	}
}

// waitingJobsHandler handles the request to get the waiting jobs.
// It encodes the waiting jobs to JSON and writes it to the response.
func waitingJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(sharedMap.Map)
	if err != nil {
		http.Error(w, "Error encoding waiting jobs to JSON", http.StatusInternalServerError)
	}
}

// finishedJobsHandler handles the request to get the finished jobs.
// It filters the jobs based on the status query parameter and encodes the result to JSON.
func finishedJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the status query parameter and convert it to lowercase
	status := r.URL.Query().Get("status")
	status = strings.ToLower(status)
	if status != "successful" && status != "failed" && status != "" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return

	}

	// Get the finished jobs based on the status
	finishedJobs, err := converter.GetRecords(status)
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
