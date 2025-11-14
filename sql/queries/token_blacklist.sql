-- name: CreateTokenBlacklist :one
INSERT INTO token_blacklist (id, user_id, token, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteExpiredBlacklistTokens :exec
DELETE FROM token_blacklist
WHERE expires_at < NOW();
