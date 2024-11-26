-- name: CreateUser :one
INSERT INTO users (id, email, password)
VALUES (
  $1,
  $2,
  $3
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
