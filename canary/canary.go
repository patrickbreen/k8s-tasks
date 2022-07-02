package main

import (
	"fmt"
	"leet/util"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func wrappedCanary(serverDomain string) {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			log.Info().Msg(fmt.Sprintf("panic occurred: %v", err))
			// increment prometheus canaryFailure counter
			canaryFailure.Add(1)
		}
	}()
	log.Info().Msg("Running canary.")
	util.RunCanary(serverDomain)
	log.Info().Msg("Ran canary.")
	// increment prometheus canarySuccess counter
	canarySuccess.Add(1)
	// record gauge with canaryDuration
	response_time := time.Since(start)
	canaryDuration.Observe(response_time.Seconds())
}

var canarySuccess = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "canary_success",
		Help: "Counter of requests.",
	},
)

var canaryFailure = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "canary_failure",
		Help: "Counter of requests.",
	},
)

var canaryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "canary_duration",
	Help: "Duration of HTTP requests.",
})

func main() {
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("Initializing canary..")

	envName := os.Getenv("ENV_NAME")
	var serverDomain = "https://tasks." + envName + ".leetcyber.com"

	// start prometheus metrics server
	prometheus.Register(canarySuccess)
	prometheus.Register(canaryFailure)
	prometheus.Register(canaryDuration)
	go http.ListenAndServe(":9000", promhttp.Handler())

	// increment prometheus canaryRuns counter
	for {
		wrappedCanary(serverDomain)
		time.Sleep(30 * time.Second)
	}

	log.Info().Msg("Canary finished successfully..")

}
