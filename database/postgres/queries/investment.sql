-- name: CreateInvestment :one
INSERT INTO investments (account_id,stock_code,company_name,total_money_buy,capital_cost,market_price,total_sell_amount,total_money_sell,current_volume,"description","status")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetAllInvestment :many
SELECT * from investments
ORDER BY stock_code;

-- name: GetInvestmentByCode :one
SELECT * from investments
where stock_code=$1;

-- name: GetInvestmentByAccountId :many
select * from investments
where account_id=$1;

-- name: UpdateInvestmentStatus :exec
update investments
set status=$2
WHERE id=$1;



