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
SELECT type, COUNT(*) AS count 
FROM reactions
WHERE chirp_id = $1
GROUP BY type 
ORDER BY count ASC;

-- name: GetUserReactions :one
SELECT type 
FROM reactions 
WHERE chirp_id = $1 AND user_id = $2;
