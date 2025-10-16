-- name: CreateOtp :one
INSERT INTO otps (id ,  otp,user_id,exp_at, created_at , updated_at) 
VALUES ($1,$2,$3,$4 ,$5,$6)

RETURNING *;

-- name: GetOtpByUserId :one
SELECT * FROM otps WHERE user_id = $1 ORDER BY exp_at DESC LIMIT 1 ;

-- -- name: DeleteOtpByUserId :one
-- SELECT * FROM otps WHERE user_id = $1;
