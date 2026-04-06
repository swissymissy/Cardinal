-- name: CreateVerificationToken :one
INSERT INTO email_verification_tokens (
    user_id, expires_at
)
VALUES (
    $1,
    NOW() + INTERVAL '30 days'
)
RETURNING *;

-- name: GetUserByVerificationToken :one
SELECT u.id, u.username, v.token, v.expires_at
FROM email_verification_tokens v 
JOIN users u ON u.id = v.user_id
WHERE v.token = $1 AND v.expires_at > NOW();

-- name: DeleteVerificationToken :exec
DELETE FROM email_verification_tokens
WHERE token = $1;

-- name: MarkUserVerified :exec
UPDATE users
SET is_verified = true 
WHERE id = $1; 

