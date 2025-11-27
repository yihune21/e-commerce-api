-- +goose Up

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL,   
    total_price NUMERIC(10,2) NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE order_items;