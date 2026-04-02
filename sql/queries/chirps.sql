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
SELECT c.*, u.username
FROM chirps c
JOIN users u ON u.id = c.user_id
ORDER BY c.created_at ASC;

-- name: GetAllChirpsFromUserID :many
SELECT c.*, u.username
FROM chirps c 
JOIN users u ON u.id = c.user_id 
WHERE c.user_id = $1
ORDER BY c.created_at ASC;

-- name: DeleteOneChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetOneChirp :one
SELECT * FROM chirps
WHERE id = $1;