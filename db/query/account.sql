-- name: CreateAccount :one
INSERT INTO accounts (
  name, owner_id, account_type_id, balance
) VALUES (
  $1, $2, $3, $4
) RETURNING *; 

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE deleted_at is null and id = $1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE deleted_at is null and id = $1
FOR NO KEY UPDATE;



-- name: GetUserAccounts :many
SELECT *
FROM accounts
WHERE deleted_at is null and owner_id = $1
order by id desc;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $2
WHERE deleted_at is null and id = $1
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE deleted_at is null and id = sqlc.arg(id)
RETURNING *;

-- name: RemoveAccount :exec
UPDATE accounts
SET deleted_at = now()
WHERE deleted_at is null and id = $1;