-- +goose Up

-- ==================== USERS TABLE INDEXES ====================
-- Email lookups for login (already has UNIQUE constraint which creates index)
-- Phone lookups for delivery
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
-- Admin filtering
CREATE INDEX idx_users_is_admin ON users(is_admin);

-- ==================== ORDERS TABLE INDEXES ====================
-- User's order history (most common query)
CREATE INDEX idx_orders_user_id ON orders(user_id);
-- Order status filtering for admins
CREATE INDEX idx_orders_order_status ON orders(order_status);
-- Payment status filtering
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
-- Delivery status for delivery personnel
CREATE INDEX idx_orders_delivery_status ON orders(delivery_status);
-- Payment method filtering
CREATE INDEX idx_orders_payment_method ON orders(payment_method);
-- Combined index for COD pending deliveries query
CREATE INDEX idx_orders_cod_pending ON orders(payment_method, delivery_status, payment_status) 
WHERE payment_method = 'cod' AND delivery_status = 'confirmed' AND payment_status = 'pending';
-- Date-based sorting and filtering
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- ==================== SHIPPING ADDRESSES TABLE INDEXES ====================
-- Lookup by order (one-to-one relationship)
CREATE UNIQUE INDEX idx_shipping_addresses_order_id ON shipping_addresses(order_id);
-- Phone number searches for delivery
CREATE INDEX idx_shipping_addresses_phone ON shipping_addresses(phone);
-- Geographic filtering
CREATE INDEX idx_shipping_addresses_city_state ON shipping_addresses(city, state);

-- ==================== ORDER ITEMS TABLE INDEXES ====================
-- Get all items for an order
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
-- Product sales analytics
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
-- Combined for order details query
CREATE INDEX idx_order_items_order_product ON order_items(order_id, product_id);

-- ==================== PRODUCTS TABLE INDEXES ====================
-- Category filtering (very common)
CREATE INDEX idx_products_category_id ON products(category_id);
-- Price range filtering
CREATE INDEX idx_products_price ON products(price);
-- Active products only
CREATE INDEX idx_products_is_active ON products(is_active) WHERE is_active = true;
-- Stock availability
CREATE INDEX idx_products_stock ON products(stock) WHERE stock > 0;
-- Product name searches (using trigram for partial matches if pg_trgm extension is available)
CREATE INDEX idx_products_name ON products(name);
-- Combined index for category + price filtering
CREATE INDEX idx_products_category_price ON products(category_id, price);

-- ==================== CATEGORIES TABLE INDEXES ====================
-- Parent-child relationships
CREATE INDEX idx_categories_parent_id ON categories(parent_id) WHERE parent_id IS NOT NULL;
-- Category name lookups
CREATE INDEX idx_categories_name ON categories(name);

-- ==================== CARTS TABLE INDEXES ====================
-- User's active cart
CREATE INDEX idx_carts_user_status ON carts(user_id, status) WHERE status = 'active';
-- All user's carts
CREATE INDEX idx_carts_user_id ON carts(user_id);

-- ==================== CART ITEMS TABLE INDEXES ====================
-- Items in a specific cart
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
-- Check if product is in cart
CREATE UNIQUE INDEX idx_cart_items_cart_product ON cart_items(cart_id, product_id);
-- Product usage analytics
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- ==================== REFRESH TOKENS TABLE INDEXES ====================
-- Token lookup for refresh
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
-- User's tokens
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- Cleanup expired tokens
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at) WHERE revoked_at IS NULL;

-- ==================== TOKEN BLACKLIST TABLE INDEXES ====================
-- Token validation check
CREATE INDEX idx_token_blacklist_token ON token_blacklist(token);
-- Cleanup expired blacklisted tokens
CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);

-- ==================== OTPs TABLE INDEXES ====================
-- OTP validation
CREATE INDEX idx_otps_user_otp ON otps(user_id, otp);
-- Cleanup expired OTPs
CREATE INDEX idx_otps_exp_at ON otps(exp_at);

-- +goose Down

-- Drop all indexes in reverse order
DROP INDEX IF EXISTS idx_otps_exp_at;
DROP INDEX IF EXISTS idx_otps_user_otp;
DROP INDEX IF EXISTS idx_token_blacklist_expires_at;
DROP INDEX IF EXISTS idx_token_blacklist_token;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_cart_items_product_id;
DROP INDEX IF EXISTS idx_cart_items_cart_product;
DROP INDEX IF EXISTS idx_cart_items_cart_id;
DROP INDEX IF EXISTS idx_carts_user_id;
DROP INDEX IF EXISTS idx_carts_user_status;
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_products_category_price;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_stock;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_order_items_order_product;
DROP INDEX IF EXISTS idx_order_items_product_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_shipping_addresses_city_state;
DROP INDEX IF EXISTS idx_shipping_addresses_phone;
DROP INDEX IF EXISTS idx_shipping_addresses_order_id;
DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_cod_pending;
DROP INDEX IF EXISTS idx_orders_payment_method;
DROP INDEX IF EXISTS idx_orders_delivery_status;
DROP INDEX IF EXISTS idx_orders_payment_status;
DROP INDEX IF EXISTS idx_orders_order_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_users_is_admin;
DROP INDEX IF EXISTS idx_users_phone;