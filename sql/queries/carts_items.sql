-- name: CreateCartItem :one
INSERT INTO cart_items (id ,  cart_id , product_id,quantity,price_at_add , created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5,$6,$7)

RETURNING *;
-- name: GetCartItemByCartIdAndProductId :one
SELECT * FROM cart_items WHERE cart_id = $1 and product_id = $2;

-- name: GetCartItemByCartId :many
SELECT * FROM cart_items WHERE cart_id = $1;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items SET quantity = $1 WHERE id = $2

RETURNING *;



