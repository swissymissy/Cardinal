-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, u.username
FROM chirps c
JOIN users u ON u.id = c.user_id
WHERE c.created_at < $1
ORDER BY c.created_at DESC
LIMIT $2;

-- name: GetAllChirpsFromUserID :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, u.username
FROM chirps c 
JOIN users u ON u.id = c.user_id 
WHERE c.user_id = $1
AND c.created_at < $2
ORDER BY c.created_at DESC
LIMIT $3;

-- name: DeleteOneChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetOneChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetFeedChirps :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, u.username
FROM chirps c 
JOIN users u ON u.id = c.user_id
WHERE (
    c.user_id = $1
    OR c.user_id IN (
        SELECT followee_id FROM followers WHERE follower_id = $1
    )
)
AND c.created_at < $2
ORDER BY c.created_at DESC
LIMIT $3;