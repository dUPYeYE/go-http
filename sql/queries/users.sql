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

-- name: UpdateUser :one
UPDATE users
SET email = $2, password = $3
WHERE id = $1
RETURNING *;

-- name: UpdateChirpyRed :exec
UPDATE users
SET is_chirpy_red = $2
WHERE id = $1;
