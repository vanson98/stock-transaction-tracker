-- name: CreateAccountEntry :one
INSERT INTO account_entries(acount_id,amount,from_type)
VALUES($1,$2,$3)
RETURNING *;

