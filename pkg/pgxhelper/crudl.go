package pgxhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = fmt.Errorf("[PGXHELPER ERROR]: %s", "not found")

func Create(ctx context.Context, querier Querier, stmt squirrel.InsertBuilder) error {
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := querier.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, querier Querier, stmt squirrel.Sqlizer) error {
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := querier.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, querier Querier, stmt squirrel.Sqlizer) error {
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := querier.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// RL is for Read List.
type RL[S any] interface {
	Total(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) (int64, error)
	GetMany(ctx context.Context, q Querier, stmt squirrel.SelectBuilder) ([]*S, error)
	GetOne(ctx context.Context, q Querier, stmt squirrel.SelectBuilder) (S, error)
	GetList(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) (*List[S], error)
}

type List[S any] struct {
	Elements []S
	Total    int64
}

func (q *rlBase[S]) GetList(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) (*List[S], error) {
	elements, err := q.GetManyObj(ctx, querier, stmt)
	if err != nil {
		return nil, err
	}

	total, err := q.Total(ctx, querier, stmt)
	if err != nil {
		return nil, err
	}

	return &List[S]{
		Elements: elements,
		Total:    total,
	}, nil
}

func (q *rlBase[S]) Total(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) (int64, error) {
	stmt = stmt.RemoveLimit().RemoveOffset()

	sql, args, err := Builder.Select("count(*)").FromSelect(stmt, "q").ToSql()
	if err != nil {
		return 0, fmt.Errorf("[PGXHELPER ERROR]: countRows %w", err)
	}

	var tmp int64
	if err = querier.QueryRow(ctx, sql, args...).Scan(&tmp); err != nil {
		err = fmt.Errorf("[PGXHELPER ERROR]: count query %w", err)
	}

	return tmp, err
}

func queryRowsPointer[S any](
	ctx context.Context,
	querier Querier,
	stmt squirrel.SelectBuilder,
) ([]*S, error) {
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slice []*S
	for rows.Next() {
		value, err := pgx.RowToStructByPos[S](rows)
		if err != nil {
			return nil, err
		}
		slice = append(slice, &value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slice, nil
}

func queryRows[S any](
	ctx context.Context,
	querier Querier,
	stmt squirrel.SelectBuilder,
) ([]S, error) {
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[S])
}

func (q *rlBase[S]) GetMany(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) ([]*S, error) {
	entities, err := queryRowsPointer[S](ctx, querier, stmt)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (q *rlBase[S]) GetManyObj(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) ([]S, error) {
	entities, err := queryRows[S](ctx, querier, stmt)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (q *rlBase[S]) GetOne(ctx context.Context, querier Querier, stmt squirrel.SelectBuilder) (S, error) {
	var defaultS S

	entities, err := q.GetMany(ctx, querier, stmt)
	if err != nil {
		return defaultS, err
	}

	if len(entities) == 0 {
		return defaultS, q.errorNotFound
	}

	return *entities[0], nil
}

type rlBase[S any] struct {
	errorNotFound error
}

func NewRL[S any](err ...error) RL[S] {
	errNotFound := ErrNotFound
	if len(err) > 0 {
		errNotFound = err[0]
	}

	return &rlBase[S]{
		errorNotFound: errNotFound,
	}
}

func SelectWhereSliceEq[T any](stmt squirrel.SelectBuilder, column string, input []T) squirrel.SelectBuilder {
	if len(input) > 0 {
		stmt = stmt.Where(squirrel.Eq{column: input})
	}
	return stmt
}

func SelectWhereEq[T comparable](stmt squirrel.SelectBuilder, column string, input T) squirrel.SelectBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.Eq{column: input})
	}
	return stmt
}

func SelectWhereNeq[T comparable](stmt squirrel.SelectBuilder, column string, input T) squirrel.SelectBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.NotEq{column: input})
	}
	return stmt
}

func SelectWhereILike[T comparable](stmt squirrel.SelectBuilder, column string, input T) squirrel.SelectBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.ILike{column: input})
	}
	return stmt
}

func WrapILike(input string) string {
	if input == "" {
		return ""
	}

	return "%" + input + "%"
}

func UpdateWhereEq[T comparable](stmt squirrel.UpdateBuilder, column string, input T) squirrel.UpdateBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.Eq{column: input})
	}
	return stmt
}

func SetMapNotNil[T any](setMap squirrel.Eq, input *T, column string) squirrel.Eq {
	if input != nil {
		setMap[column] = *input
	}
	return setMap
}

func SetMapTime(setMap squirrel.Eq, input time.Time, column string) squirrel.Eq {
	if !input.IsZero() {
		setMap[column] = input
	}
	return setMap
}

func SetMapNotEmpty[T comparable](setMap squirrel.Eq, input T, column string) squirrel.Eq {
	var defaultT T

	if input != defaultT {
		setMap[column] = input
	}

	return setMap
}

func SelectWhereGte[T comparable](stmt squirrel.SelectBuilder, column string, input T) squirrel.SelectBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.GtOrEq{column: input})
	}
	return stmt
}

func SelectWhereLte[T comparable](stmt squirrel.SelectBuilder, column string, input T) squirrel.SelectBuilder {
	var defaultT T

	if input != defaultT {
		stmt = stmt.Where(squirrel.LtOrEq{column: input})
	}
	return stmt
}
