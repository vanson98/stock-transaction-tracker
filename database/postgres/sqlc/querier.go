// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateInvestment(ctx context.Context, arg CreateInvestmentParams) (Investment, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountById(ctx context.Context, id int64) (Account, error)
	GetAllInvestment(ctx context.Context) ([]Investment, error)
	GetInvestmentByAccountId(ctx context.Context, accountID pgtype.Int8) ([]Investment, error)
	GetInvestmentByCode(ctx context.Context, stockCode string) (Investment, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error)
	UpdateInvestmentStatus(ctx context.Context, arg UpdateInvestmentStatusParams) error
}

var _ Querier = (*Queries)(nil)
