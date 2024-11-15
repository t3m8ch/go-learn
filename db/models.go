package db

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Product struct {
	Id          *uint64
	Title       string
	Description string
	Price       decimal.Decimal
}

func (p Product) String() string {
	return fmt.Sprintf("Product(id: %d, title: '%s', description: '%s', price: %s)",
		*p.Id, p.Title, p.Description, p.Price.String())
}
