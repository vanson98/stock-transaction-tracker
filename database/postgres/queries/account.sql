-- name: CreateAccount :one
INSERT INTO accounts(channel_name,"owner",balance,currency)
VALUES($1,$2,$3,$4)
RETURNING *;


-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1;

-- name: GetAccountInfoByIds :many 
SELECT a.id, a.channel_name, a.balance as cash, 
SUM(
	CASE WHEN i.capital_cost IS NULL THEN 0 ELSE (i.capital_cost * i.current_volume) END
	)
AS total_cogs,
SUM(
	CASE WHEN i.market_price IS NULL THEN 0 ELSE (i.market_price * i.current_volume) END
) AS market_value
FROM accounts as a
LEFT JOIN investments AS i ON a.id = i.account_id
WHERE a.id = ANY(@account_ids::bigint[]) AND i.status = 'active'
GROUP BY a.id,  a.channel_name, a.balance;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: ListAllAccount :many
select a.id, a.channel_name from accounts as a
where a.owner = $1;


-- name: GetAccountPaging :many
SELECT a.id, a.channel_name, a.balance, a.currency, 
SUM(
	CASE WHEN e.amount < 0 THEN e.amount ELSE 0 END
) AS withdraw, 
SUM (
	CASE WHEN e.amount > 0 THEN e.amount ELSE 0 END
) as deposit
FROM accounts AS a
LEFT JOIN entries AS e ON a.id = e.account_id
WHERE "owner" = sqlc.arg(owner)
GROUP BY a.id, a.channel_name, a.balance, a.currency;

