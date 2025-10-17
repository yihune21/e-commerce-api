-- name: CreateProduct :one
INSERT INTO products (id ,  name, description ,price,stock,category_id,image_url,is_active, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5 ,$6,$7,$8 ,$9 ,$10)

RETURNING *;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductByName :one
SELECT * FROM products WHERE name = $1;

-- -- name: UpdateProductData:one
-- UPDATE products SET password = $1 WHERE id = $2
-- RETURNING *;
-- name: DeleteProductByProductId :exec
DELETE FROM products WHERE id = $1;

