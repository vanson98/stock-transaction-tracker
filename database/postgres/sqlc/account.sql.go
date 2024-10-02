// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts(channel_name,"owner",balance,currency)
VALUES($1,$2,$3,$4)
RETURNING id, channel_name, owner, balance, currency, created_at
`

type CreateAccountParams struct {
	ChannelName string `json:"channel_name"`
	Owner       string `json:"owner"`
	Balance     int64  `json:"balance"`
	Currency    string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount,
		arg.ChannelName,
		arg.Owner,
		arg.Balance,
		arg.Currency,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAccount, id)
	return err
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, channel_name, owner, balance, currency, created_at FROM accounts
WHERE id=$1
`

func (q *Queries) GetAccountById(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, channel_name, owner, balance, currency, created_at FROM accounts
OFFSET $1 LIMIT $2
`

type ListAccountsParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, listAccounts, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.ChannelName,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING id, channel_name, owner, balance, currency, created_at
`

type UpdateAccountBalanceParams struct {
	Balance int64 `json:"balance"`
	ID      int64 `json:"id"`
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccountBalance, arg.Balance, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
