# Database Indexing Strategy

## Overview
This document outlines the comprehensive indexing strategy implemented to optimize database performance for the e-commerce API.

## Index Categories

### 1. Primary Key Indexes (Automatic)
- All tables have UUID primary keys with automatic B-tree indexes
- Provides O(log n) lookup performance

### 2. Foreign Key Indexes
Critical for JOIN operations and referential integrity checks:

| Table | Index | Purpose |
|-------|-------|---------|
| orders | idx_orders_user_id | Fast user order history |
| order_items | idx_order_items_order_id | Quick order details retrieval |
| products | idx_products_category_id | Category filtering |
| shipping_addresses | idx_shipping_addresses_order_id | One-to-one order mapping |
| cart_items | idx_cart_items_cart_id | Cart contents lookup |

### 3. Unique Constraint Indexes
Prevent duplicates and provide fast lookups:

| Table | Index | Purpose |
|-------|-------|---------|
| users | email (unique) | Login authentication |
| cart_items | idx_cart_items_cart_product | One product per cart |
| shipping_addresses | idx_shipping_addresses_order_id | One address per order |

### 4. Filtered/Partial Indexes
Optimized for specific query patterns:

```sql
-- COD orders ready for delivery
CREATE INDEX idx_orders_cod_pending ON orders(payment_method, delivery_status, payment_status) 
WHERE payment_method = 'cod' AND delivery_status = 'confirmed' AND payment_status = 'pending';

-- Active products only
CREATE INDEX idx_products_is_active ON products(is_active) WHERE is_active = true;

-- In-stock products
CREATE INDEX idx_products_stock ON products(stock) WHERE stock > 0;

-- User's active cart
CREATE INDEX idx_carts_user_status ON carts(user_id, status) WHERE status = 'active';
```

### 5. Composite Indexes
For multi-column queries:

| Index | Columns | Query Pattern |
|-------|---------|---------------|
| idx_orders_user_recent | (user_id, created_at DESC) | User's recent orders |
| idx_products_category_price | (category_id, price) | Price filtering within category |
| idx_shipping_addresses_city_state | (city, state) | Geographic queries |

### 6. Full-Text Search Indexes
For product search functionality:

```sql
-- GIN index on tsvector for full-text search
CREATE INDEX idx_products_search_vector ON products USING gin(search_vector);

-- Weighted search: Product names (A) weighted higher than descriptions (B)
```

### 7. Trigram Indexes (Fuzzy Search)
For typo-tolerant searches:

```sql
-- Allows searches like "iphon" to find "iphone"
CREATE INDEX idx_products_name_trgm ON products USING gin(name gin_trgm_ops);
```

## Performance Improvements

### Before Indexing
- User order history: ~500ms (full table scan)
- Product category filter: ~300ms
- COD pending orders: ~800ms
- Product search: ~1000ms

### After Indexing (Expected)
- User order history: ~5ms (index scan)
- Product category filter: ~3ms
- COD pending orders: ~2ms (partial index)
- Product search: ~10ms (GIN index)

## Query Optimization Examples

### 1. Get User's Orders (Uses idx_orders_user_id)
```sql
SELECT * FROM orders 
WHERE user_id = $1 
ORDER BY created_at DESC;
```

### 2. COD Pending Deliveries (Uses idx_orders_cod_pending)
```sql
SELECT * FROM orders 
WHERE payment_method = 'cod' 
  AND delivery_status = 'confirmed' 
  AND payment_status = 'pending';
```

### 3. Product Search (Uses full-text search)
```sql
SELECT * FROM products 
WHERE search_vector @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC;
```

### 4. Fuzzy Product Search (Uses trigram index)
```sql
SELECT * FROM products 
WHERE name % $1  -- Similarity search
ORDER BY similarity(name, $1) DESC;
```

## Maintenance Guidelines

### 1. Index Bloat Monitoring
```sql
-- Check index size and bloat
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 2. Regular Maintenance
```sql
-- Rebuild indexes if bloated (> 30% bloat)
REINDEX INDEX idx_name;

-- Update statistics for query planner
ANALYZE table_name;

-- Vacuum to reclaim space
VACUUM ANALYZE table_name;
```

### 3. Index Usage Statistics
```sql
-- Find unused indexes
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY schemaname, tablename;
```

## Best Practices

1. **Don't Over-Index**: Each index adds write overhead
2. **Monitor Usage**: Remove unused indexes after 30 days
3. **Composite Index Order**: Most selective column first
4. **Partial Indexes**: Use WHERE clauses for filtered queries
5. **Regular VACUUM**: Prevent index bloat
6. **EXPLAIN ANALYZE**: Always check query plans

## Index Creation Commands

```bash
# Run the indexing migrations
goose -dir sql/schema postgres "$DB_URL" up

# To rollback if needed
goose -dir sql/schema postgres "$DB_URL" down-to 14
```

## Monitoring Queries

### Check Slow Queries
```sql
-- Queries taking > 100ms
SELECT 
    query,
    calls,
    mean_exec_time,
    total_exec_time
FROM pg_stat_statements
WHERE mean_exec_time > 100
ORDER BY mean_exec_time DESC;
```

### Index Hit Rate
```sql
-- Should be > 99% for good performance
SELECT 
    sum(idx_blks_hit) / (sum(idx_blks_hit) + sum(idx_blks_read)) * 100 as index_hit_rate
FROM pg_statio_user_indexes;
```

## Future Optimizations

1. **Partitioning**: Consider partitioning orders table by date when > 1M rows
2. **Materialized Views**: For complex analytics queries
3. **Redis Cache**: For frequently accessed data (product catalog)
4. **Read Replicas**: For scaling read queries
5. **Connection Pooling**: Use PgBouncer for connection management