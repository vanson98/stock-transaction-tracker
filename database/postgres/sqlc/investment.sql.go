// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: investment.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createInvestment = `-- name: CreateInvestment :one
INSERT INTO investments (account_id,stock_code,company_name,total_money_buy,capital_cost,market_price,total_sell_amount,total_money_sell,current_volume,"description","status")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, account_id, stock_code, company_name, total_buy_amount, total_money_buy, capital_cost, market_price, total_sell_amount, total_money_sell, current_volume, description, status
`

type CreateInvestmentParams struct {
	AccountID       pgtype.Int8      `json:"account_id"`
	StockCode       string           `json:"stock_code"`
	CompanyName     pgtype.Text      `json:"company_name"`
	TotalMoneyBuy   pgtype.Numeric   `json:"total_money_buy"`
	CapitalCost     pgtype.Numeric   `json:"capital_cost"`
	MarketPrice     pgtype.Numeric   `json:"market_price"`
	TotalSellAmount int32            `json:"total_sell_amount"`
	TotalMoneySell  pgtype.Numeric   `json:"total_money_sell"`
	CurrentVolume   int32            `json:"current_volume"`
	Description     pgtype.Text      `json:"description"`
	Status          InvestmentStatus `json:"status"`
}

func (q *Queries) CreateInvestment(ctx context.Context, arg CreateInvestmentParams) (Investment, error) {
	row := q.db.QueryRow(ctx, createInvestment,
		arg.AccountID,
		arg.StockCode,
		arg.CompanyName,
		arg.TotalMoneyBuy,
		arg.CapitalCost,
		arg.MarketPrice,
		arg.TotalSellAmount,
		arg.TotalMoneySell,
		arg.CurrentVolume,
		arg.Description,
		arg.Status,
	)
	var i Investment
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.StockCode,
		&i.CompanyName,
		&i.TotalBuyAmount,
		&i.TotalMoneyBuy,
		&i.CapitalCost,
		&i.MarketPrice,
		&i.TotalSellAmount,
		&i.TotalMoneySell,
		&i.CurrentVolume,
		&i.Description,
		&i.Status,
	)
	return i, err
}

const getAllInvestment = `-- name: GetAllInvestment :many
SELECT id, account_id, stock_code, company_name, total_buy_amount, total_money_buy, capital_cost, market_price, total_sell_amount, total_money_sell, current_volume, description, status from investments
ORDER BY stock_code
`

func (q *Queries) GetAllInvestment(ctx context.Context) ([]Investment, error) {
	rows, err := q.db.Query(ctx, getAllInvestment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Investment
	for rows.Next() {
		var i Investment
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.StockCode,
			&i.CompanyName,
			&i.TotalBuyAmount,
			&i.TotalMoneyBuy,
			&i.CapitalCost,
			&i.MarketPrice,
			&i.TotalSellAmount,
			&i.TotalMoneySell,
			&i.CurrentVolume,
			&i.Description,
			&i.Status,
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

const getInvestmentByAccountId = `-- name: GetInvestmentByAccountId :many
select id, account_id, stock_code, company_name, total_buy_amount, total_money_buy, capital_cost, market_price, total_sell_amount, total_money_sell, current_volume, description, status from investments
where account_id=$1
`

func (q *Queries) GetInvestmentByAccountId(ctx context.Context, accountID pgtype.Int8) ([]Investment, error) {
	rows, err := q.db.Query(ctx, getInvestmentByAccountId, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Investment
	for rows.Next() {
		var i Investment
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.StockCode,
			&i.CompanyName,
			&i.TotalBuyAmount,
			&i.TotalMoneyBuy,
			&i.CapitalCost,
			&i.MarketPrice,
			&i.TotalSellAmount,
			&i.TotalMoneySell,
			&i.CurrentVolume,
			&i.Description,
			&i.Status,
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

const getInvestmentByCode = `-- name: GetInvestmentByCode :one
SELECT id, account_id, stock_code, company_name, total_buy_amount, total_money_buy, capital_cost, market_price, total_sell_amount, total_money_sell, current_volume, description, status from investments
where stock_code=$1
`

func (q *Queries) GetInvestmentByCode(ctx context.Context, stockCode string) (Investment, error) {
	row := q.db.QueryRow(ctx, getInvestmentByCode, stockCode)
	var i Investment
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.StockCode,
		&i.CompanyName,
		&i.TotalBuyAmount,
		&i.TotalMoneyBuy,
		&i.CapitalCost,
		&i.MarketPrice,
		&i.TotalSellAmount,
		&i.TotalMoneySell,
		&i.CurrentVolume,
		&i.Description,
		&i.Status,
	)
	return i, err
}

const updateInvestmentStatus = `-- name: UpdateInvestmentStatus :exec
update investments
set status=$2
WHERE id=$1
`

type UpdateInvestmentStatusParams struct {
	ID     int64            `json:"id"`
	Status InvestmentStatus `json:"status"`
}

func (q *Queries) UpdateInvestmentStatus(ctx context.Context, arg UpdateInvestmentStatusParams) error {
	_, err := q.db.Exec(ctx, updateInvestmentStatus, arg.ID, arg.Status)
	return err
}
