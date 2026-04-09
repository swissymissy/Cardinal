-- name: CreateNotifications :one
INSERT INTO notifications (created_at, body, receiver, triggerer, chirp_id)
VALUES (
    NOW(),
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetNotificationByReceiver :many
SELECT n.id, n.created_at, n.body, n.receiver, n.triggerer, n.chirp_id, n.is_read, u.username
FROM notifications n 
JOIN users u ON u.id = n.triggerer
WHERE n.receiver = $1
ORDER BY n.created_at DESC;

-- name: MarkOneAsRead :exec
UPDATE notifications 
SET is_read = true
WHERE id = $1 AND receiver = $2;

-- name: MarkAllAsRead :exec
UPDATE notifications
SET is_read = true 
WHERE receiver = $1;

-- name: CreateNotificationsBulk :exec
INSERT INTO notifications (body, receiver, triggerer, chirp_id, created_at)
SELECT $1, unnest($2::uuid[]), $3, $4, NOW();

-- name: CreateFollowNotification :one
INSERT INTO notifications (created_at, body, receiver, triggerer)
VALUES (NOW(), $1, $2, $3)
RETURNING *;