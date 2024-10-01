// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account_entry.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAccountEntry = `-- name: CreateAccountEntry :one
INSERT INTO account_entries(acount_id,amount,from_type)
VALUES($1,$2,$3)
RETURNING id, acount_id, amount, created_at, from_type
`

type CreateAccountEntryParams struct {
	AcountID pgtype.Int8    `json:"acount_id"`
	Amount   pgtype.Numeric `json:"amount"`
	FromType EntryFromType  `json:"from_type"`
}

func (q *Queries) CreateAccountEntry(ctx context.Context, arg CreateAccountEntryParams) (AccountEntry, error) {
	row := q.db.QueryRow(ctx, createAccountEntry, arg.AcountID, arg.Amount, arg.FromType)
	var i AccountEntry
	err := row.Scan(
		&i.ID,
		&i.AcountID,
		&i.Amount,
		&i.CreatedAt,
		&i.FromType,
	)
	return i, err
}
