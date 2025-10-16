-- name: CreateAdmin :one
INSERT INTO admins (id ,  name, email ,password, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5 ,$6)

RETURNING *;

-- name: GetAdminById :one
SELECT * FROM admins WHERE id = $1;

-- name: GetAdminByEmail :one
SELECT * FROM admins WHERE email = $1;

-- name: UpdateAdminPasword :one
UPDATE users SET password = $1 WHERE id = $2
RETURNING *;