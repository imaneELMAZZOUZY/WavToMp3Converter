package main

import (
	"github.com/gorilla/mux"

)


func routes() *mux.Router {

	mux := mux.NewRouter()
	mux.HandleFunc("/api/jobs/current", currentJobsHandler).Methods("GET")
	mux.HandleFunc("/api/jobs/finished", finishedJobsHandler).Methods("GET")
	mux.HandleFunc("/api/jobs/waiting", waitingJobsHandler).Methods("GET")
	return mux

}