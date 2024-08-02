-- name: CreateMovement :one
INSERT INTO movements (
  account_id, amount, description, validated, transfer_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetMovement :one
SELECT *
FROM movements
WHERE deleted_at is null and id = $1;

-- name: GetAccountMovements :many
SELECT *
FROM movements
WHERE deleted_at is null and account_id = $1
order by id desc
limit $2
offset $3;

-- name: UpdateMovement :one
UPDATE movements
SET amount = $2, description = $3, validated = $4
WHERE id = $1
RETURNING *;

-- name: RemoveMovement :exec
UPDATE movements
SET deleted_at = now()
WHERE deleted_at is null and id = $1;