-- name: CreateUser :exec
INSERT INTO users (
  id, username, password
) VALUES ($1, $2, $3);

-- name: UserExists :one
SELECT EXISTS (SELECT id FROM users WHERE username = $1);

