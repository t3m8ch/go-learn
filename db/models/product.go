package models

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Product struct {
	Id          *int64          `db:"id"`
	Title       string          `db:"title"`
	Description string          `db:"description"`
	Price       decimal.Decimal `db:"price"`
}

func (p Product) String() string {
	return fmt.Sprintf("Product(id: %d, title: '%s', description: '%s', price: %s)",
		*p.Id, p.Title, p.Description, p.Price.String())
}

func (p Product) TableName() string {
	return "products"
}

func (p Product) PrimaryKey() (any, any) {
	return "id", p.Id
}
