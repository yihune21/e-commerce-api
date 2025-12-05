-- name: CreateOrderItem :one
INSERT INTO order_items(id , order_id , product_id , quantity ,unit_price,total_price, created_at) 
values($1,$2,$3,$4,$5,$6 ,$7)
RETURNING *;
