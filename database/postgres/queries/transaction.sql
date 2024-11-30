-- name: GetTransactionById :one
select * from transactions
where id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions(investment_id,ticker,trading_date,trade,volume,order_price,match_volume,match_price,match_value,fee,tax,"cost","cost_of_goods_sold","return","status")
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
RETURNING *;

-- name: GetTransactionsPaging :many
SELECT T.id, to_char(T.trading_date,'dd/mm/yyyy') as trading_date, T.ticker, T.trade, T.volume, T.order_price, T.match_volume, T.match_price, T.match_value, T.fee, T.tax, T."cost", T.cost_of_goods_sold, T."return", T.status
FROM investments AS I
INNER JOIN transactions AS T ON I.id = T.investment_id
WHERE I.account_id = @account_id AND
 	  T.ticker LIKE 
	  	CASE WHEN @ticker::text = '' THEN '%%' ELSE @ticker::text END
ORDER BY trading_date DESC
OFFSET @from_offset::int LIMIT @to_limit::int;

-- name: CountTransactions :one
SELECT COUNT(T.id)
FROM investments AS I
INNER JOIN transactions AS T ON I.id = T.investment_id
WHERE I.account_id = @account_id AND
 	  T.ticker LIKE 
	  	CASE WHEN @ticker::text = '' THEN '%%' ELSE @ticker::text END;