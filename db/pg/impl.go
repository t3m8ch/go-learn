package pg

import (
	"context"
	"errors"
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
	return r.getOneSql(
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s = $1",
			strings.Join(getCols[T](), ","),
			getTableName[T](),
			key,
		),
		value,
	)
}

func (r *PgRepository[T]) GetOneSql(genSql db.GenSqlFunc, args ...any) (*T, error) {
	return r.getOneSql(genSql(getCols[T]()), args...)
}

func (r *PgRepository[T]) GetAll() ([]T, error) {
	return r.getManySql(
		fmt.Sprintf(
			"SELECT %s FROM %s",
			strings.Join(getCols[T](), ","),
			getTableName[T](),
		),
	)
}

func (r *PgRepository[T]) GetManySql(genSql db.GenSqlFunc, args ...any) ([]T, error) {
	return r.getManySql(genSql(getCols[T]()), args...)
}

func (r *PgRepository[T]) Add(entities ...T) ([]T, error) {
	sql, values, err := buildInsertSql(entities...)

	if err != nil {
		return nil, err
	}

	return r.getManySql(sql, values...)
}

func (r *PgRepository[T]) getManySql(sql string, args ...any) ([]T, error) {
	rows, err := r.pgPool.Query(context.Background(), sql, args...)
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

func (r *PgRepository[T]) getOneSql(sql string, args ...any) (*T, error) {
	rows, err := r.pgPool.Query(
		context.Background(),
		sql,
		args...,
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

func buildInsertSql[T db.Entity](entities ...T) (string, []any, error) {
	cols := getCols[T]()
	returningSql := fmt.Sprintf("RETURNING %s", strings.Join(cols, ","))

	pkIdx, err := getPkIdx(cols, getPrimaryKey[T]())
	if err != nil {
		return "", nil, err
	}

	cols = append(cols[:pkIdx], cols[pkIdx+1:]...)

	var sqlSb strings.Builder
	sqlSb.WriteString(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES\n",
		getTableName[T](),
		strings.Join(cols, ","),
	))

	values := pushValuesToSql(&sqlSb, entities, cols, pkIdx)
	sqlSb.WriteString(returningSql)

	return sqlSb.String(), values, nil
}

func pushValuesToSql[T db.Entity](sqlSb *strings.Builder, entities []T, cols []string, pkIdx int) []any {
	values := make([]any, 0, len(cols)*len(entities))

	for i, entity := range entities {
		sqlSb.WriteString("(")
		entityValues := getValues(entity)
		entityValues = append(entityValues[:pkIdx], entityValues[pkIdx+1:]...)

		for j := range entityValues {
			num := i*len(entityValues) + j + 1
			sqlSb.WriteString(fmt.Sprintf("$%d", num))

			if j < len(entityValues)-1 {
				sqlSb.WriteString(", ")
			}
		}

		sqlSb.WriteString(")")

		if i != len(entities)-1 {
			sqlSb.WriteString(",")
		}

		sqlSb.WriteString("\n")
		values = append(values, entityValues...)
	}

	return values
}

func getPkIdx(cols []string, pk string) (int, error) {
	for i, v := range cols {
		if v == pk {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("PK '%s' not found", pk))
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

func getValues[T db.Entity](entity T) []any {
	value := reflect.ValueOf(entity)
	result := make([]any, 0, value.NumField())

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		result = append(result, fieldValue.Interface())
	}

	return result
}

func getPrimaryKey[T db.Entity]() string {
	var e T
	pk, _ := e.PrimaryKey()
	return pk
}
