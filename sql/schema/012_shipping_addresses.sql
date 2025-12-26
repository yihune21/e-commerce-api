-- +goose Up

CREATE TABLE shipping_addresses (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    address_line1 TEXT NOT NULL,
    address_line2 TEXT,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    country TEXT NOT NULL DEFAULT 'Ethiopia',
    delivery_instructions TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE shipping_addresses;