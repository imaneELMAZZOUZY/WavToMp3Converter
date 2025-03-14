package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (appDep *appDep) routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/jobs/current", appDep.currentJobsHandler)
	mux.HandleFunc("GET /api/jobs/finished", appDep.finishedJobsHandler)
	mux.HandleFunc("GET /api/jobs/waiting", appDep.waitingJobsHandler)
	mux.Handle("GET /api/metrics", promhttp.Handler())
	return mux

}
