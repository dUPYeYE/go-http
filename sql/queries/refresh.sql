-- name: AddRefreshToken :one
INSERT INTO refreshtokens (user_id, token)
VALUES (
  $1,
  $2
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refreshtokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refreshtokens
SET revoked_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC', expires_at = null
WHERE token = $1;
