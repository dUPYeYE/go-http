-- name: AddChirp :one
INSERT INTO chirps (id, user_id, body)
VALUES (
  $1,
  $2,
  $3
)
RETURNING *;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirpsFromUser :many
SELECT * FROM chirps WHERE user_id = $1;

-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: RemoveChirp :exec
DELETE FROM chirps WHERE id = $1;
