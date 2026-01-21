package pgxhelper

import (
	"context"
	"os"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Value[T any](value *T) T {
	var defaultValue T

	if value == nil {
		return defaultValue
	}
	return *value
}

type TransactionClosure func(ctx context.Context, tx pgx.Tx) error

func InTransaction(ctx context.Context, pg *pgxpool.Pool, closure TransactionClosure) (err error) {
	tx, err := pg.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(ctx)
	}()

	if err := closure(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func ToAny(ids []string) []any {
	tmp := make([]interface{}, len(ids))
	for i, val := range ids {
		tmp[i] = val
	}
	return tmp
}

func ApplyFixture(ctx context.Context, pg *pgxpool.Pool, fixture func() string) error {
	return PerformSQL(ctx, pg, fixture())
}

func PerformSQL(ctx context.Context, pg *pgxpool.Pool, query string, args ...interface{}) error {
	err := InTransaction(ctx, pg, func(ctx context.Context, tx pgx.Tx) error {
		for _, q := range strings.Split(query, ";") {
			if _, err := tx.Exec(ctx, q, args...); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
