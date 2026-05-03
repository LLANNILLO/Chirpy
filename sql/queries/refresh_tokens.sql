-- name: CreateRefresh :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES(
  $1,
  $2,
  $3
)
RETURNING *;

-- name: GetRefresh :one
SELECT * FROM refresh_tokens
WHERE token = $1
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
AND refresh_tokens.revoked_at IS NULL
AND refresh_tokens.expires_at > NOW()
LIMIT 1;
