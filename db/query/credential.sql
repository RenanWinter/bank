-- name: CreateCredential :one
INSERT INTO credentials (
  user_id, password
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserActiveCredential :one
SELECT *
FROM credentials
WHERE deleted_at is null and user_id = $1
order by id desc
limit 1;

-- name: GetUserCredentials :many
SELECT *
FROM credentials
WHERE user_id = $1
order by id desc;

-- name: RemoveUserCredential :exec
UPDATE credentials
SET deleted_at = now()
WHERE deleted_at is null and user_id = $1;