-- name: CreateOrder :one
INSERT INTO orders (id , user_id , order_status , total ,payment_status, created_at , updated_at) 
values($1,$2,$3,$4,$5,$6 , $7)
RETURNING *;
