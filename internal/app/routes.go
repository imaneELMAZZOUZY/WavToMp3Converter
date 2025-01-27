package app

import (
	"net/http"
	
)


func (appDep *AppDep) Routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/jobs/current", appDep.currentJobsHandler)
	mux.HandleFunc("GET /api/jobs/finished", appDep.finishedJobsHandler)
	mux.HandleFunc("GET /api/jobs/waiting", appDep.waitingJobsHandler)
	return mux

}