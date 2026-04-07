-- +goose Up
CREATE TABLE reactions (
    chirp_id        UUID NOT NULL, 
    user_id         UUID NOT NULL, 
    type            TEXT NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, chirp_id),
    FOREIGN KEY (chirp_id) REFERENCES chirps(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)  REFERENCES users(id)  ON DELETE CASCADE,
    CONSTRAINT valid_reaction CHECK (type IN ('❤️', '😂', '😮', '😢', '👍'))
);

CREATE INDEX idx_reactions_chirp_id ON reactions(chirp_id);

-- +goose Down
DROP TABLE reactions;