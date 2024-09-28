-- name: CreateAccount :one
insert into accounts(channel_name,"owner",balance,buy_fee,sell_free,currency)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: ListAccounts :many
select * from accounts
order by "owner";

-- name: UpdateAccount :exec
UPDATE accounts
SET channel_name = $1
WHERE id = $2;



