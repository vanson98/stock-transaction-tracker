-- name: CreateUser :one 
INSERT INTO users(username,email)
VALUES ($1,$2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

