-- +goose Up

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

ALTER TABLE products
ADD COLUMN IF NOT EXISTS search_vector tsvector;

UPDATE products
SET search_vector =
    setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B');

CREATE INDEX IF NOT EXISTS idx_products_search_vector
ON products USING gin(search_vector);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION products_search_vector_update()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER products_search_vector_trigger
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION products_search_vector_update();
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_products_name_trgm
ON products USING gin(name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_categories_name_trgm
ON categories USING gin(name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_name_trgm
ON users USING gin(name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_users_email_trgm
ON users USING gin(email gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_products_active_category_price
ON products(category_id, price)
WHERE is_active = true AND stock > 0;

CREATE INDEX IF NOT EXISTS idx_orders_user_recent
ON orders(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_orders_status_recent
ON orders(order_status, created_at DESC);


-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS products_search_vector_trigger ON products;
DROP FUNCTION IF EXISTS products_search_vector_update();
-- +goose StatementEnd

ALTER TABLE products
DROP COLUMN IF EXISTS search_vector;

DROP INDEX IF EXISTS idx_orders_status_recent;
DROP INDEX IF EXISTS idx_orders_user_recent;
DROP INDEX IF EXISTS idx_products_active_category_price;

DROP INDEX IF EXISTS idx_users_email_trgm;
DROP INDEX IF EXISTS idx_users_name_trgm;
DROP INDEX IF EXISTS idx_categories_name_trgm;
DROP INDEX IF EXISTS idx_products_name_trgm;

DROP INDEX IF EXISTS idx_products_search_vector;
