-- name: CreateAccountEntry :one
INSERT INTO account_entries(account_id,amount,from_type)
VALUES($1,$2,$3)
RETURNING *;

-- name: GetAccountEntryById :one
SELECT * FROM account_entries
WHERE id=$1;

