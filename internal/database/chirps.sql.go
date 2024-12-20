// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addChirp = `-- name: AddChirp :one
INSERT INTO chirps (id, user_id, body)
VALUES (
  $1,
  $2,
  $3
)
RETURNING id, user_id, body, created_at, updated_at
`

type AddChirpParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Body   string
}

func (q *Queries) AddChirp(ctx context.Context, arg AddChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, addChirp, arg.ID, arg.UserID, arg.Body)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, user_id, body, created_at, updated_at
FROM chirps
ORDER BY created_at ASC
`

func (q *Queries) GetAllChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirp = `-- name: GetChirp :one
SELECT id, user_id, body, created_at, updated_at FROM chirps WHERE id = $1
`

func (q *Queries) GetChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getChirpsFromUser = `-- name: GetChirpsFromUser :many
SELECT id, user_id, body, created_at, updated_at FROM chirps WHERE user_id = $1 ORDER BY created_at ASC
`

func (q *Queries) GetChirpsFromUser(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirpsFromUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeChirp = `-- name: RemoveChirp :exec
DELETE FROM chirps WHERE id = $1
`

func (q *Queries) RemoveChirp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeChirp, id)
	return err
}
