-- name: CreateOrder :one
INSERT INTO orders (id , user_id , order_status , total ,payment_status, payment_method, delivery_status, created_at , updated_at) 
values($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetAllOrders :many
SELECT * FROM orders;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUserId :many
SELECT * FROM orders WHERE user_id = $1;

-- name: GetOrdersPerPage :many
SELECT * FROM orders ORDER BY created_at Limit $1;


-- name: UpdateOrderStatus :one
UPDATE orders SET order_status = $1, updated_at = $3 WHERE id = $2
RETURNING *;

-- name: UpdateDeliveryStatus :one
UPDATE orders SET delivery_status = $1, updated_at = $3 WHERE id = $2
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE orders SET payment_status = $1, updated_at = $3 WHERE id = $2
RETURNING *;

-- name: CompletePaymentOnDelivery :one
UPDATE orders 
SET payment_status = 'completed', 
    order_status = 'delivered',
    delivery_status = 'delivered',
    updated_at = $2
WHERE id = $1 AND payment_method = 'cod'
RETURNING *;