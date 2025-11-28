-- name: CreateOrder :one
INSERT INTO orders (id , user_id , status , total , created_at , updated_at) 
values($1,$2,$3,$4,$5,$6)
RETURNING *;
