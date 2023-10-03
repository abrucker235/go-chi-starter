package state

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type State struct {
	server *http.Server
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *State {
	router := chi.NewRouter()

	if viper.GetBool("metrics.enabled") {
		router.Get("/metrics", promhttp.Handler().ServeHTTP)
	}

	//TODO: Add liveness and readiness endpoints

	return &State{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", viper.GetString("state.port")),
			Handler: router,
		},
		logger: logger.With().Str("sub-system", "state").Logger(),
	}
}

func (s *State) Serve(signal context.Context, timeout context.Context) {
	shutdown := make(chan error)
	go func(signal context.Context, timeout context.Context, shutdown chan error) {
		<-signal.Done()
		shutdown <- s.server.Shutdown(timeout)
	}(signal, timeout, shutdown)

	if err := s.server.ListenAndServe(); err != nil || err != http.ErrServerClosed {
		shutdown <- err
	}

	if err := <-shutdown; err != nil {
		s.logger.Err(err)
	}
}
