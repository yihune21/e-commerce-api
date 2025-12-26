-- name: SearchProducts :many
-- Full-text search on products
SELECT * FROM products 
WHERE search_vector @@ plainto_tsquery('english', $1)
  AND is_active = true
  AND stock > 0
ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC
LIMIT $2;

-- name: FuzzySearchProducts :many
-- Fuzzy search for typo tolerance (e.g., "iphon" finds "iphone")
SELECT * FROM products 
WHERE name % $1  -- Similarity operator
  AND is_active = true
  AND stock > 0
ORDER BY similarity(name, $1) DESC
LIMIT $2;

-- name: SearchProductsByPriceRange :many
-- Optimized price range search using index
SELECT * FROM products
WHERE price >= $1 
  AND price <= $2
  AND is_active = true
  AND stock > 0
ORDER BY price ASC
LIMIT $3;

-- name: SearchProductsByCategoryAndPrice :many
-- Uses composite index on category_id and price
SELECT * FROM products
WHERE category_id = $1
  AND price >= $2
  AND price <= $3
  AND is_active = true
  AND stock > 0
ORDER BY price ASC;

-- name: GetPopularProducts :many
-- Get most ordered products (uses order_items index)
SELECT 
    p.*,
    COUNT(oi.id) as order_count,
    SUM(oi.quantity) as total_quantity_sold
FROM products p
INNER JOIN order_items oi ON p.id = oi.product_id
WHERE p.is_active = true
GROUP BY p.id
ORDER BY order_count DESC
LIMIT $1;

-- name: SearchUsersByName :many
-- Fuzzy search users by name (for admin panel)
SELECT id, name, email, phone, is_admin, created_at
FROM users
WHERE name % $1
ORDER BY similarity(name, $1) DESC
LIMIT $2;

-- name: GetCODPendingDeliveries :many
-- Uses the partial index idx_orders_cod_pending
SELECT o.*, sa.* 
FROM orders o
INNER JOIN shipping_addresses sa ON o.id = sa.order_id
WHERE o.payment_method = 'cod' 
  AND o.delivery_status = 'confirmed' 
  AND o.payment_status = 'pending'
ORDER BY o.created_at ASC;

-- name: GetUserRecentOrders :many
-- Uses composite index idx_orders_user_recent
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetOrdersByGeography :many
-- Uses idx_shipping_addresses_city_state
SELECT o.*, sa.*
FROM orders o
INNER JOIN shipping_addresses sa ON o.id = sa.order_id
WHERE sa.city = $1 
  AND sa.state = $2
ORDER BY o.created_at DESC;

-- name: GetLowStockProducts :many
-- Uses partial index on stock
SELECT * FROM products
WHERE stock > 0 AND stock <= $1
  AND is_active = true
ORDER BY stock ASC;

-- name: AutocompleteProductName :many
-- Fast product name autocomplete using trigram index
SELECT DISTINCT name FROM products
WHERE name ILIKE $1 || '%'
  AND is_active = true
ORDER BY name
LIMIT 10;