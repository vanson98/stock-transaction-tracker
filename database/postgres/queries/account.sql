-- name: ListAccounts :many
select * from accounts
order by "owner";

