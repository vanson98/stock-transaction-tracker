-- name: CreateAccount :one
INSERT INTO accounts(channel_name,"owner",balance,currency)
VALUES($1,$2,$3,$4)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
OFFSET $1 LIMIT $2;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1;
