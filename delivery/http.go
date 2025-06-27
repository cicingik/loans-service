package delivery

import (
	"fmt"
	"net/http"

	"github.com/cicingik/loans-service/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	// HTTPEngine ...
	HTTPEngine struct {
		Mux      *chi.Mux
		HTTPHost string
		HTTPPort int
	}
)

// NewHTTPServer ...
func NewHTTPServer(cfg *config.AppConfig) *HTTPEngine {
	return &HTTPEngine{
		Mux:      chi.NewMux(),
		HTTPHost: cfg.HTTPHost,
		HTTPPort: cfg.HTTPPort,
	}
}

// InitMiddleware ...
func (h *HTTPEngine) InitMiddleware(appMiddleware ...func(http.Handler) http.Handler) {
	c := h.Mux

	//  rateLimit := customMiddleware.RateLimit(1*time.Second, 3)

	c.Use(middleware.RequestID)
	c.Use(middleware.RealIP)
	c.Use(middleware.Logger)
	c.Use(middleware.Recoverer)

	// App-level middleware
	for _, m := range appMiddleware {
		c.Use(m)
	}

	c.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	))
}

// RegisterHandler ...
func (h *HTTPEngine) RegisterHandler(registerFn func(*chi.Mux)) {
	registerFn(h.Mux)
}

// Serve ...
func (h *HTTPEngine) Serve() error {
	binding := fmt.Sprintf("%s:%d", h.HTTPHost, h.HTTPPort)
	fmt.Printf("Running HTTP Server in %s \n", binding)
	return http.ListenAndServe(binding, h.Mux)
}
