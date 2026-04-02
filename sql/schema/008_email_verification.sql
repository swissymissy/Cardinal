-- +goose Up
ALTER TABLE users ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE email_verification_tokens (
    token       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE email_verification_tokens;
ALTER TABLE users DROP COLUMN is_verified;

