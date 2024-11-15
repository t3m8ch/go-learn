package db

type ProductRepository interface {
	GetById(id uint64) (*Product, error)
	GetAll() (*[]Product, error)
	Save(product *Product) (*Product, error)
}
