-- name: CreateUser :exec
INSERT INTO users (
  id, username, password
) VALUES ($1, $2, $3);

