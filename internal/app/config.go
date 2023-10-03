package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abrucker235/go-chi-starter/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	cmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type App struct {
	server *http.Server
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *App {
	app := chi.NewRouter()
	app.Use(cmiddleware.RequestID)
	app.Use(cmiddleware.RealIP)
	app.Use(cmiddleware.Logger)
	app.Use(cmiddleware.Recoverer)

	if viper.GetBool("metrics.enabled") {
		app.Use(middleware.PrometheusMiddleware)
	}

	//TODO: Adding router
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", viper.GetString("http.port")),
			Handler: app,
		},
		logger: logger.With().Str("sub-system", "app-server").Logger(),
	}
}

func (a *App) Serve(serve context.Context, timeout context.Context) {
	shutdown := make(chan error)
	go func(serve context.Context, timeout context.Context, shutdown chan error) {
		<-serve.Done()
		shutdown <- a.server.Shutdown(timeout)
	}(serve, timeout, shutdown)

	if err := a.server.ListenAndServe(); err != nil || err != http.ErrServerClosed {
		shutdown <- err
	}

	if err := <-shutdown; err != nil {
		a.logger.Err(err)
	}
}
