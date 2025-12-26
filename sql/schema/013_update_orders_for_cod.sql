-- +goose Up

ALTER TABLE orders 
ADD COLUMN payment_method VARCHAR(20) DEFAULT 'cod' CHECK (payment_method IN ('cod', 'online', 'card')),
ADD COLUMN delivery_status VARCHAR(30) DEFAULT 'pending' CHECK (delivery_status IN ('pending', 'confirmed', 'out_for_delivery', 'delivered', 'failed'));

-- Update existing orders to have COD payment method
UPDATE orders SET payment_method = 'cod' WHERE payment_method IS NULL;

-- +goose Down
ALTER TABLE orders 
DROP COLUMN payment_method,
DROP COLUMN delivery_status;