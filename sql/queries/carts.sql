-- name: CreateCart :one
INSERT INTO carts (id ,  user_id, status , created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5 )

RETURNING *;

-- name: GetCartByUserId :one
SELECT * FROM carts WHERE user_id = $1;


