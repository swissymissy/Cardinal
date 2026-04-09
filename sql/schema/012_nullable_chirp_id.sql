-- +goose Up
ALTER TABLE notifications 
ALTER COLUMN chirp_id 
DROP NOT NULL;

-- +goose Down
ALTER TABLE notifications
ALTER COLUMN chirp_id
SET NOT NULL;
