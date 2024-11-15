package main

import (
	"context"
	"fmt"
	"go-learn/cli"
	"go-learn/db"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	pool, err := pgxpool.New(
		context.Background(),
		"postgres://t3m8ch@localhost/productsdb",
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer pool.Close()

	db.Init(pool)
	productRepo := db.CreatePgProductRepository(pool)
	cli.Loop(&productRepo)
}
