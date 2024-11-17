package db

import "errors"

type pkey = string
type pvalue = any

type Entity interface {
	TableName() string
	PrimaryKey() (pkey, pvalue)
}

var ErrNotFound = errors.New("Not found")

type GenSqlFunc = func(cols []string) string

type Repository[T Entity] interface {
	GetOne(key any, value any) (*T, error)
	GetOneSql(genSql GenSqlFunc, args ...any) (*T, error)
	GetAll() ([]T, error)
	GetManySql(genSql GenSqlFunc, args ...any) ([]T, error)
	Add(entities ...T) ([]T, error)
	// Update(entity *T) error
	// DeleteOne(key any, value any) error
}
