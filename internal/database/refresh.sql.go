// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addRefreshToken = `-- name: AddRefreshToken :one
INSERT INTO refreshtokens (user_id, token)
VALUES (
  $1,
  $2
)
RETURNING token, user_id, created_at, updated_at, expires_at, revoked_at
`

type AddRefreshTokenParams struct {
	UserID uuid.UUID
	Token  string
}

func (q *Queries) AddRefreshToken(ctx context.Context, arg AddRefreshTokenParams) (Refreshtoken, error) {
	row := q.db.QueryRowContext(ctx, addRefreshToken, arg.UserID, arg.Token)
	var i Refreshtoken
	err := row.Scan(
		&i.Token,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT token, user_id, created_at, updated_at, expires_at, revoked_at FROM refreshtokens WHERE token = $1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (Refreshtoken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshToken, token)
	var i Refreshtoken
	err := row.Scan(
		&i.Token,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refreshtokens
SET revoked_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC', updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC', expires_at = null
WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, token)
	return err
}
