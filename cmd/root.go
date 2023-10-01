package cmd

import (
	"github.com/abrucker235/go-chi-starter/cmd/migrate"
	"github.com/abrucker235/go-chi-starter/cmd/server"
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{}

func init() {
	RootCMD.AddCommand(migrate.MigrateCMD)
	RootCMD.AddCommand(server.ServerCMD)
}
