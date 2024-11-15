package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("Product not found")

type PgProductRepository struct {
	pgPool *pgxpool.Pool
}

func CreatePgProductRepository(pgPool *pgxpool.Pool) PgProductRepository {
	return PgProductRepository{pgPool}
}

func (repo *PgProductRepository) GetById(id uint64) (*Product, error) {
	var p Product

	row := repo.pgPool.QueryRow(
		context.Background(),
		"SELECT id, title, description, price FROM products WHERE id = $1",
		id,
	)
	err := row.Scan(&p.Id, &p.Title, &p.Description, &p.Price)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (repo *PgProductRepository) GetAll() (*[]Product, error) {
	rows, err := repo.pgPool.Query(
		context.Background(),
		"SELECT id, title, description, price FROM products",
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.Title, &p.Description, &p.Price)
		if err != nil {
			return &products, err
		}
		products = append(products, p)
	}

	return &products, nil
}

func (repo *PgProductRepository) Save(product Product) (*Product, error) {
	var query string
	var args pgx.NamedArgs

	if product.Id != nil {
		query = `
			INSERT INTO products (id, title, description, price)
			VALUES (@id, @title, @description, @price)
			ON CONFLICT (id) DO UPDATE
			SET title = @title, description = @description, price = @price
			RETURNING id
		`
		args = pgx.NamedArgs{
			"id":          product.Id,
			"title":       product.Title,
			"description": product.Description,
			"price":       product.Price,
		}
	} else {
		query = `
			INSERT INTO products (title, description, price)
			VALUES (@title, @description, @price)
			RETURNING id
		`
		args = pgx.NamedArgs{
			"title":       product.Title,
			"description": product.Description,
			"price":       product.Price,
		}
	}

	var id uint64
	err := repo.pgPool.QueryRow(context.Background(), query, args).Scan(&id)

	if err != nil {
		return nil, err
	}

	if product.Id == nil {
		product.Id = &id
	}

	return &product, nil
}
