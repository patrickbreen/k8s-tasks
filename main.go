package main

import (
	"leet/util"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("Initializing server..")

	// init postgres and prometheus
	defer util.DBFree()
	util.InitPostgres()
	InitPrometheus()
	mux := InitAppServer()

	log.Info().Msg("Server initialized")
	go http.ListenAndServe(":9000", promhttp.Handler())
	http.ListenAndServe(":8080", mux)
}
