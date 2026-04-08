-- name: CreateComment :one
INSERT INTO comments (chirp_id, user_id, body)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCommentsByChirpID :many
SELECT c.id, c.chirp_id, c.user_id, c.body, c.created_at, c.updated_at, u.username
FROM comments c 
JOIN users u ON u.id = c.user_id
WHERE c.chirp_id = $1
ORDER BY c.created_at ASC;

-- name: DeleteComment :one
DELETE FROM comments 
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: EditComment :one
UPDATE comments 
SET body = $1, 
    updated_at = NOW()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: GetCommentCount :one
SELECT COUNT(*) AS count 
FROM comments
WHERE chirp_id = $1;
