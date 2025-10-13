-- name: CreateOtp :one
INSERT INTO otps (id ,  otp,user_id,exp_at, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5,$6)

RETURNING *;

-- name: GetOtpByUserId :one
SELECT * FROM users WHERE id = $1;
