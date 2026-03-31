-- +goose UP
ALTER TABLE users 
ADD COLUMN username TEXT NOT NULL DEFAULT '';

UPDATE users
SET username = id::text
WHERE username = '';

ALTER TABLE users
ADD CONSTRAINT users_username_key UNIQUE (username);

-- +goose Down
ALTER TABLE users 
DROP COLUMN username;
