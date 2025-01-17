package main

import (
	"net/http"
)


func routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/jobs/current", currentJobsHandler)
	mux.HandleFunc("GET /api/jobs/finished", finishedJobsHandler)
	mux.HandleFunc("GET /api/jobs/waiting", waitingJobsHandler)
	return mux

}