-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE deleted_at is null and id = $1;

-- name: RemoveTransfer :exec
UPDATE transfers
SET deleted_at = now()
WHERE deleted_at is null and id = $1;