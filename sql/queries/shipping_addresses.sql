-- name: CreateShippingAddress :one
INSERT INTO shipping_addresses(
    id,
    order_id,
    full_name,
    phone,
    address_line1,
    address_line2,
    city,
    state,
    postal_code,
    country,
    delivery_instructions,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetShippingAddressByOrderId :one
SELECT * FROM shipping_addresses
WHERE order_id = $1;

-- name: UpdateShippingAddress :one
UPDATE shipping_addresses
SET 
    full_name = $2,
    phone = $3,
    address_line1 = $4,
    address_line2 = $5,
    city = $6,
    state = $7,
    postal_code = $8,
    country = $9,
    delivery_instructions = $10,
    updated_at = $11
WHERE id = $1
RETURNING *;

-- name: DeleteShippingAddress :exec
DELETE FROM shipping_addresses
WHERE id = $1;