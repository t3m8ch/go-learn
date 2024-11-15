package main

import (
	"context"
	"fmt"
	"go-learn/cli"
	"go-learn/db"
	"log"
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

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS products (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2)
		);
	`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	productRepo := db.CreatePgProductRepository(pool)
	cli.Loop(&productRepo)
}
