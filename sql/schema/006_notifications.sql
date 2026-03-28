-- +goose Up
CREATE TABLE notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    body        TEXT NOT NULL,
    receiver    UUID NOT NULL,
    triggerer   UUID NOT NULL,
    chirp_id    UUID NOT NULL,
    is_read     BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (receiver) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (triggerer) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (chirp_id) REFERENCES chirps(id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_receiver ON notifications(receiver);

-- +goose Down
DROP TABLE notifications;