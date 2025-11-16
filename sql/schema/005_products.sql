
-- +goose Up

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    image_url TEXT, 
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose Down
DROP TABLE products;
