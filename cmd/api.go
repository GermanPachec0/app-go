package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/GermanPachec0/app-go/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 4000
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			db, err := NewDatabasePool(ctx, 16)
			if err != nil {
				return err
			}
			defer db.Close()

			api := api.NewAPI(ctx, db)
			srv := api.Server(port)

			go func() { _ = srv.ListenAndServe() }()

			slog.Info("started api", "port", port)

			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			return nil
		},
	}

	return cmd
}

func NewDatabasePool(ctx context.Context, maxConns int) (*pgxpool.Pool, error) {
	if maxConns == 0 {
		maxConns = 1
	}
	url := fmt.Sprintf(
		"%s?pool_max_conns=%d&pool_min_conns=%d",
		os.Getenv("DATABASE_CONNECTION_POOL_URL"),
		maxConns,
		2,
	)

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	// Setting the build statement cache to nil helps this work with pgbouncer
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second
	return pgxpool.NewWithConfig(ctx, config)
}
