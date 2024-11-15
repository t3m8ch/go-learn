package main

import (
	"bufio"
	"context"
	"fmt"
	"go-learn/db"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		var cmd string
		fmt.Scan(&cmd)

		switch cmd {
		case "add":
			var product db.Product

			fmt.Print("title: ")
			scanner.Scan()
			product.Title = scanner.Text()

			fmt.Print("description: ")
			scanner.Scan()
			product.Description = scanner.Text()

			fmt.Print("price: ")
			scanner.Scan()
			priceStr := scanner.Text()

			price, err := decimal.NewFromString(priceStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid price")
				continue
			}
			product.Price = price

			p, err := productRepo.Save(product)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(p)
		case "getall":
			products, err := productRepo.GetAll()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, p := range *products {
				fmt.Println(p)
			}
		case "get":
			fmt.Print("id: ")
			scanner.Scan()
			id, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid id")
				continue
			}

			product, err := productRepo.GetById(id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(product)
		}
	}
}
