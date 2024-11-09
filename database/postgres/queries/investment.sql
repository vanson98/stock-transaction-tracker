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

-- name: GetInvestmentsByAccountId :many
select * from investments
where account_id=$1;

-- name: GetInvestmentById :one
select * from investments
where id=$1;

-- name: UpdateInvestmentStatus :exec
update investments
set status=$2
WHERE id=$1;

-- name: UpdateInvestmentWhenBuying :exec
update investments
set buy_volume = $2,
buy_value = $3,
capital_cost = $4,
current_volume = $5,
fee = $6,
tax = $7,
updated_date = $8
where id = $1;



