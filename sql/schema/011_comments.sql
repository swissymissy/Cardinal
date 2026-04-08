-- +goose Up
CREATE TABLE comments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chirp_id    UUID NOT NULL,
    user_id     UUID NOT NULL,
    body        TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (chirp_id) REFERENCES chirps(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)  REFERENCES users(id)  ON DELETE CASCADE
);

CREATE INDEX idx_comments_chirp_id ON comments(chirp_id);

-- +goose Down
DROP TABLE comments;