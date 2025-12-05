-- name: CreateOrder :one
INSERT INTO orders (id , user_id , order_status , total ,payment_status, created_at , updated_at) 
values($1,$2,$3,$4,$5,$6 , $7)
RETURNING *;

-- name: GetAllOrders :many
SELECT * FROM orders;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUserId :many
SELECT * FROM orders WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders SET order_status = $1 WHERE id = $2
RETURNING *;