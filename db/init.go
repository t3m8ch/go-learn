package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(pgPool *pgxpool.Pool) {
	_, err := pgPool.Exec(context.Background(), `
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
}
