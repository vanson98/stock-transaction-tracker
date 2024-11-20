-- name: CreateInvestment :one
insert into investments(account_id,ticker,company_name,buy_volume,buy_value,capital_cost,market_price,sell_volume,sell_value,current_volume,description,status,fee,tax)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: SearchInvestmentPaging :many
SELECT * from investments
WHERE account_id=sqlc.arg(account_id) AND (ticker ILIKE @search_text::text OR company_name ILIKE @search_text::text)
ORDER BY 
    CASE WHEN @order_by::text = 'ticker' AND @sort_type::text = 'ascending' THEN ticker END ASC,
    CASE WHEN @order_by::text = 'ticker' AND @sort_type::text = 'descending' THEN ticker END DESC,
    CASE WHEN @order_by::text = 'status' AND @sort_type::text = 'ascending' THEN "status" END ASC,
    CASE WHEN @order_by::text = 'status' AND @sort_type::text = 'descending' THEN "status" END DESC
OFFSET @from_offset::int LIMIT @take_limit::int;

-- name: CountInvestment :one
SELECT COUNT(*) from investments
WHERE account_id=sqlc.arg(account_id) AND (ticker ILIKE @search_text::text OR company_name ILIKE @search_text::text);

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

-- name: UpdateInvestmentWhenSeling :exec
UPDATE investments
SET sell_volume = sell_volume + @sell_transaction_volume,
sell_value = sell_value + @sell_transaction_value,
current_volume = current_volume - @sell_transaction_volume,
fee = fee + @transaction_fee,
tax = tax + @transaction_tax, 
updated_date = sqlc.arg(updated_date)
WHERE id = $1;



