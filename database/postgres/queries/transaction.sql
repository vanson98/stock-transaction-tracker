-- name: CreateTransaction :one
INSERT INTO transactions(investment_id,ticker,trading_date,trade,volume,order_price,match_volume,match_price,match_value,fee,tax,"cost","cost_of_goods_sold","return","status")
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
RETURNING *;

-- name: UpdateTransactionCost :exec
update transactions
set cost = $2,
cost_of_goods_sold = $3
where id = $1;
