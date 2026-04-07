-- name: AddReaction :one
INSERT INTO reactions (chirp_id, user_id, type)
VALUES ($1, $2, $3)
ON CONFLICT (chirp_id, user_id)
DO UPDATE SET type = EXCLUDED.type 
RETURNING *;

-- name: RemoveReaction :exec
DELETE FROM reactions
WHERE chirp_id = $1 AND user_id = $2;

-- name: GetReactionsByChirpID :many
SELECT r.chirp_id, r.user_id, r.type, r.created_at, u.username
FROM reactions r 
JOIN users u ON u.id = r.user_id
WHERE r.chirp_id = $1
ORDER BY r.created_at DESC;

-- name: GetReactionCounts :many
SELECT type, COUNT(*) as count
FROM reactions
WHERE chirp_id = $1
GROUP BY type 
ORDER BY count ASC;

-- name: GetUserReactions :one
SELECT type 
FROM reactions 
WHERE chirp_id = $1 AND user_id = $2;
