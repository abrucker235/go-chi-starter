package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/abrucker235/go-chi-starter/internal/app"
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

	ServerCMD.Flags().String("metrics-port", "", "Metrics Port")
	viper.BindPFlag("metrics.port", ServerCMD.Flags().Lookup("metrics-port"))
	viper.BindEnv("metrics.port", "METRICS-PORT")
	viper.SetDefault("metrics.port", "9100")
}

func server(cmd *cobra.Command, args []string) {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	serve, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	application := app.NewApp(logger)

	timeout, cancel := context.WithTimeout(serve, viper.GetDuration("shutdown.gracePeriod"))
	defer cancel()

	go application.Serve(serve, timeout)

	//TODO: need to start up listner for internal traffic like liveness, readiness, health, metrics
	if viper.GetBool("metrics.enabled") {
		go func() {

		}()
	}

	<-serve.Done()
	//TODO: gracefully shut things down
	<-timeout.Done()
}
