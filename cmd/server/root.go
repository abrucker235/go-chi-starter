package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/abrucker235/go-chi-starter/internal/app"
	"github.com/abrucker235/go-chi-starter/internal/state"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ServerCMD = &cobra.Command{
	Use:   "server",
	Short: "start server",
	Run:   server,
}

func init() {
	ServerCMD.Flags().String("http-port", "", "HTTP Port")
	viper.BindPFlag("http.port", ServerCMD.Flags().Lookup("http-port"))
	viper.BindEnv("http.port", "HTTP-PORT")
	viper.SetDefault("http.port", "8080")

	ServerCMD.Flags().Bool("metrics-enabled", false, "Enable Prometheus Metrics")
	viper.BindPFlag("metrics.enabled", ServerCMD.Flags().Lookup("metrics-enabled"))
	viper.BindEnv("metrics.enabled", "METRICS-ENABLED")
	viper.SetDefault("metrics.enabled", "false")

	ServerCMD.Flags().String("state-port", "", "Application State Port (ready, liveness, metrics)")
	viper.BindPFlag("state.port", ServerCMD.Flags().Lookup("state-port"))
	viper.BindEnv("state.port", "METRICS-PORT")
	viper.SetDefault("state.port", "9080")
}

func server(cmd *cobra.Command, args []string) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	serve, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	application := app.New(logger)
	state := state.New(logger)

	timeout, cancel := context.WithTimeout(serve, viper.GetDuration("shutdown.gracePeriod"))
	defer cancel()

	go application.Serve(serve, timeout)
	go state.Serve(serve, timeout)

	<-serve.Done()
	//TODO: gracefully shut things down
	<-timeout.Done()
}
