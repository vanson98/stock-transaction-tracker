package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IStore interface {
	Querier
	ExecTx(ctx context.Context, fn func(q *Queries) (interface{}, error)) (interface{}, error)
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

func (store *DBStore) ExecTx(ctx context.Context, fn func(q *Queries) (interface{}, error)) (interface{}, error) {
	tx, err := store.connectionPool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	result, err := fn(store.Queries)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	} else {
		cmErr := tx.Commit(ctx)
		if cmErr != nil {
			err = fmt.Errorf("tx commit err : %v", cmErr)
		}
	}

	return result, err
}

// func (store *DBStore) ExecTx(ctx context.Context, fn interface{}) (interface{}, error) {
// 	tx, err := store.connectionPool.BeginTx(ctx, pgx.TxOptions{
// 		IsoLevel: pgx.ReadCommitted,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	fnValue := reflect.ValueOf(fn)
// 	params := []reflect.Value{reflect.ValueOf(store.Queries)}
// 	results := fnValue.Call(params)

// 	if len(results) != 2 {
// 		return nil, fmt.Errorf("fn must return (T, error)")
// 	}

// 	if !results[1].IsNil() {
// 		err = results[1].Interface().(error)
// 	}

// 	if err != nil {
// 		if rbErr := tx.Rollback(ctx); rbErr != nil {
// 			err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
// 		}
// 	} else {
// 		cmErr := tx.Commit(ctx)
// 		if cmErr != nil {
// 			err = fmt.Errorf("tx commit err : %v", cmErr)
// 		}
// 	}

// 	return results[0].Interface(), err
// }
