package store

import (
	"context"
	"strings"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func InitializeCQL(ctx context.Context, logger *zerolog.Logger) *gocql.Session {
	cluster := gocql.NewCluster(viper.GetString("db.host"))
	cluster.Keyspace = viper.GetString("db.keyspace")

	cluster.Compressor = &gocql.SnappyCompressor{}

	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: viper.GetInt("db.retries")}

	switch strings.ToLower(viper.GetString("db.consistency")) {
	case "any":
		cluster.Consistency = gocql.Any
	case "one":
		cluster.Consistency = gocql.One
	case "two":
		cluster.Consistency = gocql.Two
	case "three":
		cluster.Consistency = gocql.Three
	case "quorum":
		cluster.Consistency = gocql.Quorum
	case "all":
		cluster.Consistency = gocql.All
	case "localone":
		cluster.Consistency = gocql.LocalOne
	case "eachquorum":
		cluster.Consistency = gocql.EachQuorum
	default:
		cluster.Consistency = gocql.LocalQuorum
	}

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatal().Err(err)
	}

	return session
}
