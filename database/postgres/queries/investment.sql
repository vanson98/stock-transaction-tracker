-- name: CreateInvestment :one
insert into investments(account_id,ticker,company_name,buy_volume,buy_value,capital_cost,market_price,sell_volume,sell_value,current_volume,description,status,fee,tax)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetInvestmentById :one
select * from investments
where id=$1;

-- name: SearchInvestmentPaging :many
SELECT * FROM investment_overview
WHERE account_id = ANY(@account_ids::bigint[]) 
    AND ticker ILIKE 
        CASE WHEN @search_text::text = '' THEN '%%' ELSE '%' || @search_text::text || '%' END
ORDER BY 
    CASE WHEN @order_by::text = 'ticker' AND @sort_type::text = 'ascending' THEN ticker END ASC,
    CASE WHEN @order_by::text = 'ticker' AND @sort_type::text = 'descending' THEN ticker END DESC,
    CASE WHEN @order_by::text = 'status' AND @sort_type::text = 'ascending' THEN "status" END ASC,
    CASE WHEN @order_by::text = 'status' AND @sort_type::text = 'descending' THEN "status" END DESC,
    CASE WHEN @order_by::text = 'channel_name' AND @sort_type::text = 'descending' THEN "channel_name" END DESC,
    CASE WHEN @order_by::text = 'channel_name' AND @sort_type::text = 'ascending' THEN "channel_name" END ASC,
    CASE WHEN @order_by::text = 'profit' AND @sort_type::text = 'ascending' THEN profit END ASC,
    CASE WHEN @order_by::text = 'profit' AND @sort_type::text = 'descending' THEN profit END DESC
OFFSET @from_offset::int LIMIT @take_limit::int;

-- name: CountInvestment :one
SELECT COUNT(*) from investments
WHERE account_id=ANY(@account_ids::bigint[]) AND (ticker ILIKE @search_text::text OR company_name ILIKE @search_text::text);

-- name: GetInvestmentByTicker :one
SELECT * from investments
where ticker=$1 AND account_id =$2;

-- name: UpdateInvestmentWhenBuying :one
update investments
set buy_volume = buy_volume + @buy_transaction_volume,
buy_value = buy_value + @buy_transaction_value,
capital_cost = @capital_cost,
current_volume = current_volume + @buy_transaction_volume,
fee = fee + @transaction_fee,
tax = tax + @transaction_tax,
updated_date = @updated_date, 
status = @status
where id = $1
RETURNING *;

-- name: UpdateInvestmentWhenSeling :one
UPDATE investments
SET sell_volume = sell_volume + @sell_transaction_volume,
sell_value = sell_value + @sell_transaction_value,
current_volume = current_volume - @sell_transaction_volume,
fee = fee + @transaction_fee,
tax = tax + @transaction_tax, 
status= @status,
updated_date = sqlc.arg(updated_date)
WHERE id = $1
RETURNING *;


