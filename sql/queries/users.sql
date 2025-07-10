-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
    VALUES (gen_random_uuid (), now(), now(), $1, $2)
RETURNING
    *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: GetUserFromRefreshToken :one
SELECT
    *
FROM
    users
WHERE
    users.id = (
        SELECT
            refresh_tokens.user_id
        FROM
            refresh_tokens
        WHERE
            refresh_tokens.token = $1);

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $2,
    hashed_password = $3,
    updated_at = timezone('UTC', now())
WHERE
    id = $1
RETURNING
    email;

