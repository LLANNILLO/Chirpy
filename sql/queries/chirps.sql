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

-- name: ListChirps :many
SELECT id, body, created_at, updated_at, user_id FROM chirps
  ORDER BY created_at ASC;

-- name: ListChirpsByAuthor :many
SELECT id, body, created_at, updated_at, user_id FROM chirps
  WHERE user_id = $1
  ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT id, body, created_at, updated_at, user_id
  FROM chirps
  WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
  WHERE user_id = $1
  AND id = $2;
