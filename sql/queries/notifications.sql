-- name: CreateNotifications :one
INSERT INTO notifications (created_at, body, receiver, triggerer, chirp_id)
VALUES (
    NOW(),
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetNotificationByReceiver :many
SELECT * FROM notifications
WHERE receiver = $1
ORDER BY created_at DESC;

-- name: MarkOneAsRead :exec
UPDATE notifications 
SET is_read = true
WHERE id = $1 AND receiver = $2;

-- name: MarkAllAsRead :exec
UPDATE notifications
SET is_read = true 
WHERE receiver = $1;

