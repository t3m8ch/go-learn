package main

import (
	"context"
	"fmt"
	"go-learn/cli"
	"go-learn/db/models"
	"go-learn/db/pg"
	"os"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := pgxpool.ParseConfig("postgres://t3m8ch@localhost/productsdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse connection string: %v\n", err)
		os.Exit(1)
	}
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer pool.Close()

	pg.Init(pool)
	productRepo := pg.CreatePgRepository[models.Product](pool)
	cli.Loop(&productRepo)
}
