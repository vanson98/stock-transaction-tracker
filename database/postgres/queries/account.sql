-- name: CreateAccount :one
INSERT INTO accounts(channel_name,"owner",balance,buy_fee,sell_free,currency)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
OFFSET $1 LIMIT $2;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id=$1;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1;
