-- name: CreateInvestment :one
insert into investments(account_id,ticker,company_name,buy_volume,buy_value,capital_cost,market_price,sell_volume,sell_value,current_volume,description,status,fee,tax)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetAllInvestment :many
SELECT * from investments
ORDER BY ticker;

-- name: GetInvestmentByTicker :one
SELECT * from investments
where ticker=$1;

-- name: GetInvestmentByAccountId :many
select * from investments
where account_id=$1;

-- name: UpdateInvestmentStatus :exec
update investments
set status=$2
WHERE id=$1;



