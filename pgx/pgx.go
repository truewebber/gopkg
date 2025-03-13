package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoAtomic(
	ctx context.Context,
	pool *pgxpool.Pool,
	handle func(context.Context, pgx.Tx) error,
) error {
	if err := DoAtomicWithOptions(ctx, pool, pgx.TxOptions{}, handle); err != nil {
		return fmt.Errorf("do atomic with options: %w", err)
	}

	return nil
}

func DoAtomicWithOptions(
	ctx context.Context,
	pool *pgxpool.Pool,
	txOptions pgx.TxOptions,
	handle func(context.Context, pgx.Tx) error,
) error {
	tx, err := pool.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err = handle(ctx, tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("%v: rollback tx: %w", err, rbErr)
		}

		return fmt.Errorf("do atomic operation on tx: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
