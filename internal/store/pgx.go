package store

import (
	"context"
	"fmt"

	"github.com/abrucker235/go-chi-starter/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func InitializePostgres(ctx context.Context, logger *zerolog.Logger) *pgxpool.Pool {
	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s pool_min_conns=%s pool_max_conns=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.name"),
		viper.GetString("db.min-connections"),
		viper.GetString("db.max-connections"),
	)

	if schema := viper.GetString("db.schema"); schema != "" {
		connection = fmt.Sprintf("%s search_path=%s", connection, schema)
	}

	if password := viper.GetString("db.password"); password != "" {
		connection = fmt.Sprintf("%s password=%s", connection, password)
	}

	sslMode := "sslmode=disable"
	if caFile := viper.GetString("db.ca-cert"); caFile != "" && utils.FileExists(caFile) {
		sslMode = fmt.Sprintf("sslmode=verify-full sslrootcert=%s", caFile)

		if clientCert := viper.GetString("db.client-cert"); clientCert != "" && utils.FileExists(clientCert) {
			sslMode = fmt.Sprintf("%s sslcert=%s", sslMode, clientCert)
		}

		if clientKey := viper.GetString("db.client-key"); clientKey != "" && utils.FileExists(clientKey) {
			sslMode = fmt.Sprintf("%s sslkey=%s", sslMode, clientKey)
		}
	}

	connection = fmt.Sprintf("%s %s", connection, sslMode)

	configuration, err := pgxpool.ParseConfig(connection)
	if err != nil {
		logger.Fatal().Err(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, configuration)
	if err != nil {
		logger.Fatal().Err(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal().Err(err)
	}

	return pool
}
