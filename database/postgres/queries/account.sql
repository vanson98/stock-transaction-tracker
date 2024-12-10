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

-- name: GetAccountInfoById :one 
select a.id, a.channel_name, a.balance,a.currency, a."owner",
SUM(
	case
	WHEN amount > 0 and e.type='TM' then amount
	ELSE 0
	END
	) as deposit,
SUM(
	CASE 
	WHEN amount < 0 and e.type='TM' THEN amount
	ELSE 0 
	END
) AS withdrawal
from accounts as a
left join entries as e on a.id = e.account_id 
where a.id = $1
GROUP BY a.id,  a.channel_name, a.balance, a.currency, a."owner"
LIMIT 1;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1;

-- name: ListAllAccount :many
select a.id, a.channel_name from accounts as a;


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

