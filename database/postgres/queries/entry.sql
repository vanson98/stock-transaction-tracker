-- name: CreateEntry :one
INSERT INTO entries(account_id,amount,"type")
VALUES($1,$2,$3)
RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id=$1;

