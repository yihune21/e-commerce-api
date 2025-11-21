-- name: CreateCategoty :one
INSERT INTO categories (id ,  name, description ,parent_id, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5 ,$6)
RETURNING *;


-- name: GetCategoryById :one
SELECT * FROM categories WHERE id = $1;

-- name: UpdateCategoryName :one
UPDATE categories SET name = $1 WHERE id = $2
RETURNING *;
