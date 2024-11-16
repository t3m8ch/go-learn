package pg

import (
	"context"
	"fmt"
	"go-learn/db"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository[T db.Entity] struct {
	pgPool *pgxpool.Pool
}

func CreatePgRepository[T db.Entity](pgPool *pgxpool.Pool) PgRepository[T] {
	return PgRepository[T]{pgPool}
}

func (r *PgRepository[T]) GetOne(key any, value any) (*T, error) {
	rows, err := r.pgPool.Query(
		context.Background(),
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s = $1",
			strings.Join(getCols[T](), ","),
			getTableName[T](),
			key,
		),
		value,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, db.ErrNotFound
	}

	values, err := rows.Values()

	if err != nil {
		return nil, err
	}

	entity, err := fillEntity[T](values)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *PgRepository[T]) GetAll() ([]T, error) {
	rows, err := r.pgPool.Query(
		context.Background(),
		fmt.Sprintf(
			"SELECT %s FROM %s",
			strings.Join(getCols[T](), ","),
			getTableName[T](),
		),
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var entities []T
	for rows.Next() {
		values, err := rows.Values()

		if err != nil {
			return entities, err
		}

		entity, err := fillEntity[T](values)

		if err != nil {
			return entities, err
		}

		entities = append(entities, *entity)
	}

	return entities, nil
}

func getCols[T db.Entity]() []string {
	var entity T

	v := reflect.ValueOf(entity)
	cols := make([]string, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		if col := fieldType.Tag.Get("db"); col != "" {
			cols = append(cols, col)
		}
	}

	return cols
}

func getTableName[T db.Entity]() string {
	var entity T
	return entity.TableName()
}

func fillEntity[T db.Entity](valuesFromDB []any) (*T, error) {
	var e T
	entityType := reflect.TypeOf(e)

	entityValue := reflect.New(entityType).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityValue.Field(i)

		if field.Type().Kind() == reflect.Ptr {
			newValue := reflect.New(field.Type().Elem())
			newValue.Elem().Set(reflect.ValueOf(valuesFromDB[i]))
			field.Set(newValue)
		} else {
			field.Set(reflect.ValueOf(valuesFromDB[i]))
		}
	}

	entity := entityValue.Interface().(T)
	return &entity, nil
}
