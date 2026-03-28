-- +goose Up
CREATE TABLE chirps (
    id              UUID PRIMARY KEY NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    body            TEXT NOT NULL,  
    user_id         UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_chirps_user_id ON chirps(user_id);

-- +goose Down
DROP TABLE chirps;