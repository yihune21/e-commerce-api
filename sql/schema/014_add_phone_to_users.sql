-- +goose Up

ALTER TABLE users 
ADD COLUMN phone TEXT;

-- +goose Down
ALTER TABLE users 
DROP COLUMN phone;