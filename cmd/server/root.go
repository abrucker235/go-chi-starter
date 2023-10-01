package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"

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
	viper.SetDefault("http.port", "3000")

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	app_listener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("http.port")))
	if err != nil {
	}

	var metrics_listener net.Listener
	if viper.GetBool("metrics.enabled") {
		metrics_listener, err = net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("metrics.port")))
		if err != nil {
		}
		http.Serve(metrics_listener, nil)
	}

	http.Serve(app_listener, nil)

	<-ctx.Done()
}
