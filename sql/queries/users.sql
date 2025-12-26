-- name: CreateUser :one
INSERT INTO users (id ,  name, email ,password, phone, is_admin, created_at , updated_at) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)

RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserPasword :one
UPDATE users SET password = $1 WHERE id = $2
RETURNING *;

-- name: UpdateUserPhone :one
UPDATE users SET phone = $1, updated_at = $3 WHERE id = $2
RETURNING *;

-- name: DeleteUserByUserId :exec
DELETE FROM users WHERE id = $1;
