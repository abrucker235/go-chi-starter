package cmd

import (
	"strings"

	"github.com/abrucker235/go-chi-starter/cmd/migrate"
	"github.com/abrucker235/go-chi-starter/cmd/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCMD = &cobra.Command{
	Use:   "go-chi-starter",
	Short: "Go Chi Starter App",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch strings.ToLower(viper.GetString("log.level")) {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	},
}

func init() {
	RootCMD.PersistentFlags().String("log-level", "", "log level (debug, info, warn, error, fatal, panic)")
	viper.BindPFlag("log.level", RootCMD.Flags().Lookup("log-level"))
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.SetDefault("log.level", "info")

	RootCMD.AddCommand(migrate.MigrateCMD)
	RootCMD.AddCommand(server.ServerCMD)
}
