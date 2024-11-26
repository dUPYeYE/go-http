-- +goose Up
CREATE TABLE refreshtokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '60 days',
    revoked_at TIMESTAMP,

    UNIQUE (user_id)
);

-- +goose Down
DROP TABLE refreshtokens;
