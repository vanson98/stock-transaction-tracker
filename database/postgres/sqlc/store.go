package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IStore interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

type DBStore struct {
	connectionPool *pgxpool.Pool
	*Queries
}

func NewStore(cnnPool *pgxpool.Pool) IStore {
	return &DBStore{
		connectionPool: cnnPool,
		Queries:        New(cnnPool),
	}
}

func (store *DBStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connectionPool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	err = fn(store.Queries)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
