-- name: FollowUser :one
INSERT INTO followers (follower_id, followee_id)
VALUES ($1, $2)
ON CONFLICT (follower_id, followee_id) DO NOTHING
RETURNING *;

-- name: UnfollowUser :exec
DELETE FROM followers 
WHERE follower_id = $1 AND followee_id = $2;

-- name: GetFollowings :many
SELECT f.followee_id, f.created_at, u.username
FROM followers f
JOIN users u ON u.id = f.followee_id 
WHERE f.follower_id = $1 
ORDER BY f.created_at DESC;

-- name: GetFollowers :many
SELECT f.follower_id, f.created_at, u.username 
FROM followers f
JOIN users u ON u.id = f.followers_id
WHERE f.followee_id = $1
ORDER BY f.created_at DESC;

-- name: GetCountFollowings :one 
SELECT COUNT(*) AS following_count
FROM followers
WHERE follower_id = $1;

-- name: GetCountFollowers :one
SELECT COUNT(*) AS follower_count
FROM followers 
WHERE followee_id = $1;

-- name: GetFollowersEmail :many
SELECT u.id , u.email
FROM followers f
JOIN users u ON u.id = f.follower_id 
WHERE f.followee_id = $1
ORDER BY f.created_at DESC;
