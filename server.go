package main

import (
	"fmt"
	"leet/tasks"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

func health(w http.ResponseWriter, req *http.Request) {
	log.Info().Msg(fmt.Sprintf("log health check"))
	fmt.Fprintf(w, "healthy\n")
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "my_http_requests_total",
		Help: "Counter of requests.",
	},
	[]string{"path", "method"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "my_http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path", "method"})

var inFlightRequests = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "my_http_requests_in_flight",
		Help: "Gauge of in flight requests.",
	},
	[]string{"path"},
)

type WrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{w, http.StatusOK}
}

func (wrw *WrappedResponseWriter) WriteHeader(code int) {
	wrw.StatusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}

type MyWrappedHandler struct {
	Handler http.Handler
}

// handle logging and metrics
func (m *MyWrappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// monitoring
	start := time.Now()
	inFlightRequests.WithLabelValues(r.URL.Path).Inc()
	wrw := NewWrappedResponseWriter(w)

	// call handler
	m.Handler.ServeHTTP(wrw, r)

	// monitoring
	status_code := fmt.Sprintf("%d", wrw.StatusCode)
	totalRequests.WithLabelValues(r.URL.Path, r.Method).Add(1)
	response_time := time.Since(start)
	httpDuration.WithLabelValues(r.URL.Path, r.Method).Observe(response_time.Seconds())
	inFlightRequests.WithLabelValues(r.URL.Path).Dec()

	//logging
	// TODO, add any panic stack traces with line numbers
	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("response_time", response_time.String()).
		Str("source_ip", r.RemoteAddr).
		Str("referer", r.Referer()).
		Str("response_code", status_code).Msg("")
}

func NewMyWrappedHandler(handlerToWrap http.Handler) *MyWrappedHandler {
	return &MyWrappedHandler{handlerToWrap}
}

func InitPrometheus() {
	prometheus.Register(totalRequests)
	prometheus.Register(httpDuration)
	prometheus.Register(inFlightRequests)
}

func InitAppServer() *MyWrappedHandler {
	// register http routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	fileserver := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fileserver))
	mux.HandleFunc("/api/v1/tasks/", tasks.TasksHandler)
	// middleware
	return NewMyWrappedHandler(mux)
}
