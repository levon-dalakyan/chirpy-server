-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
    VALUES ($1, now(), now(), $2, now() + interval '60 days')
RETURNING
    *;

-- name: GetRefreshToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token = $1;

-- name: RevokeRefreshToken :exec
UPDATE
    refresh_tokens
SET
    revoked_at = timezone('UTC', now()),
    updated_at = timezone('UTC', now())
WHERE
    token = $1;

