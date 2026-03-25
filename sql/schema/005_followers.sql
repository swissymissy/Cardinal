-- +goose up
CREATE TABLE followers (
    follower_id     UUID NOT NULL,
    followee_id     UUID NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT no_self_follow CHECK (follower_id <> followee_id)
);

CREATE INDEX idx_followers_followee_id ON followers(followee_id);

-- +goose down
DROP TABLE followers;