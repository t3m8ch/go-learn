package db

import "errors"

type pkey = any
type pvalue = any

type Entity interface {
	TableName() string
	PrimaryKey() (pkey, pvalue)
}

var ErrNotFound = errors.New("Not found")
var ErrConflict = errors.New("Conflict")

type Repository[T Entity] interface {
	GetOne(key any, value any) (*T, error)
	GetOneSql(genSql func(cols []string) string, args ...any) (*T, error)
	GetAll() ([]T, error)
	GetManySql(genSql func(cols []string) string, args ...any) ([]T, error)
	// ExecuteSql(sql string, args ...any) (any, error)
	// Add(entities ...T) (*T, error)
	// AddIgnoreCoflict(entities ...*T) (*T, error)
	// Update(entity *T) error
	// DeleteOne(key any, value any) error
}
