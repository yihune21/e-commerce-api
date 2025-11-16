-- name: CreateCategoty :one
INSERT INTO categories (id ,  name, description ,parent_id, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5 ,$6)
RETURNING *;