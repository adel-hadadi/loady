package transport

import (
	"fmt"
	"github.com/adel-hadadi/load-balancer/internal/lb"
	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Server struct {
	controller *lb.Controller
}

func New(controller *lb.Controller) *Server {
	return &Server{
		controller: controller,
	}
}

func (s *Server) Serve(port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.routes(),
	}

	return srv.ListenAndServe()
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "A histogram of latencies for requests.",
	}, []string{"server", "path", "method", "code"})

	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_total_operations_processed",
		Help: "The total number of processed requests",
	}, []string{"server", "path", "method", "code"})

	registry := prometheus.NewRegistry()
	registry.MustRegister(requestDuration, opsProcessed, lb.HealthyServers)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sv, err := s.controller.Next(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sv == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		metrics := httpsnoop.CaptureMetrics(sv.Serve(w, r), w, r)

		opsProcessed.With(prometheus.Labels{
			"server": sv.Name(),
			"path":   r.URL.Path,
			"method": r.Method,
			"code":   fmt.Sprintf("%d", metrics.Code),
		}).Inc()

		requestDuration.With(prometheus.Labels{
			"server": sv.Name(),
			"path":   r.URL.Path,
			"method": r.Method,
			"code":   fmt.Sprintf("%d", metrics.Code),
		}).Observe(metrics.Duration.Seconds())
	})

	mux.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	return mux
}
