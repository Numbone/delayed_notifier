-- name: GetUserByID :one
SELECT id, email, name, created_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name)
VALUES ($1, $2)
RETURNING id, email, name, created_at;

-- name: ListUsers :many
SELECT id, email, name, created_at
FROM users;