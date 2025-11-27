-- +goose Up

CREATE TABLE orders (
    id UUID PRIMARY KEY ,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total NUMERIC(10,2) NOT NULL,    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE orders;