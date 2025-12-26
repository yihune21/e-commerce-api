-- +goose Up

-- Enable PostgreSQL extensions for advanced text search
-- Note: You may need superuser privileges to create extensions
-- If these fail, ask your DBA to enable them
CREATE EXTENSION IF NOT EXISTS pg_trgm; -- For fuzzy text matching
CREATE EXTENSION IF NOT EXISTS unaccent; -- For accent-insensitive search

-- ==================== FULL TEXT SEARCH INDEXES ====================

-- Products full-text search
ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Update search vector for existing products
UPDATE products SET search_vector = 
    setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B');

-- Create GIN index for full-text search
CREATE INDEX idx_products_search_vector ON products USING gin(search_vector);

-- Trigger to automatically update search vector
CREATE OR REPLACE FUNCTION products_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER products_search_vector_trigger 
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION products_search_vector_update();

-- ==================== TRIGRAM INDEXES FOR FUZZY SEARCH ====================

-- Fuzzy product name search (typo-tolerant)
CREATE INDEX idx_products_name_trgm ON products USING gin(name gin_trgm_ops);

-- Fuzzy category name search
CREATE INDEX idx_categories_name_trgm ON categories USING gin(name gin_trgm_ops);

-- Fuzzy user name search (for admin panel)
CREATE INDEX idx_users_name_trgm ON users USING gin(name gin_trgm_ops);

-- Fuzzy email search (for admin panel)
CREATE INDEX idx_users_email_trgm ON users USING gin(email gin_trgm_ops);

-- ==================== COMPOSITE INDEXES FOR COMPLEX QUERIES ====================

-- Products: Active products sorted by price within category
CREATE INDEX idx_products_active_category_price 
ON products(category_id, price) 
WHERE is_active = true AND stock > 0;

-- Orders: Recent orders by user
CREATE INDEX idx_orders_user_recent 
ON orders(user_id, created_at DESC);

-- Orders: Admin dashboard - recent orders by status
CREATE INDEX idx_orders_status_recent 
ON orders(order_status, created_at DESC);

-- +goose Down

-- Drop triggers and functions
DROP TRIGGER IF EXISTS products_search_vector_trigger ON products;
DROP FUNCTION IF EXISTS products_search_vector_update();

-- Drop search vector column
ALTER TABLE products DROP COLUMN IF EXISTS search_vector;

-- Drop composite indexes
DROP INDEX IF EXISTS idx_orders_status_recent;
DROP INDEX IF EXISTS idx_orders_user_recent;
DROP INDEX IF EXISTS idx_products_active_category_price;

-- Drop trigram indexes
DROP INDEX IF EXISTS idx_users_email_trgm;
DROP INDEX IF EXISTS idx_users_name_trgm;
DROP INDEX IF EXISTS idx_categories_name_trgm;
DROP INDEX IF EXISTS idx_products_name_trgm;

-- Drop full-text search indexes
DROP INDEX IF EXISTS idx_products_search_vector;

-- Note: We don't drop extensions as they might be used by other schemas