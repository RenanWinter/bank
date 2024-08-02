-- name: CreateUser :one
INSERT INTO users (
  username, email, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * 
FROM users 
WHERE deleted_at is null and username = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE deleted_at is null and id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE deleted_at is null and email = $1;

-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at is null
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2, name = $3, email = $4, updated_at = now()
WHERE deleted_at is null and id = $1
RETURNING *;

-- name: RemoveUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;