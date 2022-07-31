package main

import (
	"crypto/tls"
	"fmt"
	"leet/util"
	"net/http"
	"os"
	"runtime/debug"
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
			log.Info().Msg(fmt.Sprintf("panic occurred: %v, stacktrace: %s", err, string(debug.Stack())))
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

	// bad security, but I'm doing this in this toy project to get through the self-signed certs
	cert, _ := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
		},
	}
	http.DefaultClient = client

	envName := os.Getenv("ENV_NAME")
	var serverDomain = "https://tasks." + envName + ".leetcyber.com"
	log.Info().Msg("serverDomain is: " + serverDomain)

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
