-- Database Performance Monitoring Queries
-- Run these periodically to monitor your database health

-- ============================================
-- 1. INDEX USAGE AND EFFECTIVENESS
-- ============================================

-- Check which indexes are being used
SELECT 
    schemaname AS schema,
    tablename AS table,
    indexname AS index,
    idx_scan AS times_used,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
    CASE 
        WHEN idx_scan = 0 THEN 'UNUSED - Consider dropping'
        WHEN idx_scan < 100 THEN 'RARELY USED'
        WHEN idx_scan < 1000 THEN 'OCCASIONALLY USED'
        ELSE 'FREQUENTLY USED'
    END AS usage_category
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- ============================================
-- 2. SLOW QUERIES (Requires pg_stat_statements)
-- ============================================

-- Find slowest queries
-- Note: Requires pg_stat_statements extension
-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

/*
SELECT 
    substring(query, 1, 100) AS query_preview,
    calls,
    round(mean_exec_time::numeric, 2) AS avg_time_ms,
    round(total_exec_time::numeric, 2) AS total_time_ms,
    round(100.0 * total_exec_time / sum(total_exec_time) OVER (), 2) AS percentage_of_total_time
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY mean_exec_time DESC
LIMIT 20;
*/

-- ============================================
-- 3. TABLE SIZES AND BLOAT
-- ============================================

-- Check table and index sizes
SELECT 
    schemaname AS schema,
    tablename AS table,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- ============================================
-- 4. CACHE HIT RATES
-- ============================================

-- Overall cache hit rate (should be > 99%)
SELECT 
    'Overall' AS category,
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) * 100 AS cache_hit_rate
FROM pg_statio_user_tables
UNION ALL
SELECT 
    'Indexes' AS category,
    sum(idx_blks_hit) / (sum(idx_blks_hit) + sum(idx_blks_read)) * 100 AS cache_hit_rate
FROM pg_statio_user_indexes;

-- Per-table cache hit rates
SELECT 
    schemaname,
    tablename,
    heap_blks_hit AS cache_hits,
    heap_blks_read AS disk_reads,
    CASE 
        WHEN heap_blks_hit + heap_blks_read = 0 THEN 0
        ELSE round(100.0 * heap_blks_hit / (heap_blks_hit + heap_blks_read), 2)
    END AS cache_hit_rate
FROM pg_statio_user_tables
ORDER BY heap_blks_hit + heap_blks_read DESC;

-- ============================================
-- 5. LOCK MONITORING
-- ============================================

-- Check for blocking locks
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocked_activity.usename AS blocked_user,
    blocking_locks.pid AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocked_activity.query AS blocked_statement,
    blocking_activity.query AS current_statement_in_blocking_process
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks 
    ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;

-- ============================================
-- 6. CONNECTION MONITORING
-- ============================================

-- Current connections by state
SELECT 
    state,
    count(*) AS connections,
    max(query_start) AS oldest_query_start
FROM pg_stat_activity
WHERE pid <> pg_backend_pid()
GROUP BY state
ORDER BY connections DESC;

-- ============================================
-- 7. VACUUM AND ANALYZE STATUS
-- ============================================

-- Tables that need VACUUM or ANALYZE
SELECT 
    schemaname,
    tablename,
    n_live_tup AS live_tuples,
    n_dead_tup AS dead_tuples,
    round(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_tuple_percent,
    last_vacuum,
    last_autovacuum,
    last_analyze,
    last_autoanalyze,
    CASE 
        WHEN n_dead_tup > 1000 AND round(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) > 20 
        THEN 'NEEDS VACUUM'
        ELSE 'OK'
    END AS vacuum_status
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;

-- ============================================
-- 8. SEQUENTIAL SCAN MONITORING
-- ============================================

-- Tables with high sequential scan rate (might need indexes)
SELECT 
    schemaname,
    tablename,
    seq_scan,
    idx_scan,
    CASE 
        WHEN seq_scan + idx_scan = 0 THEN 0
        ELSE round(100.0 * seq_scan / (seq_scan + idx_scan), 2)
    END AS seq_scan_percentage,
    n_live_tup AS table_rows,
    CASE 
        WHEN seq_scan > idx_scan AND n_live_tup > 10000 THEN 'NEEDS INDEX'
        WHEN seq_scan > idx_scan AND n_live_tup > 1000 THEN 'CONSIDER INDEX'
        ELSE 'OK'
    END AS recommendation
FROM pg_stat_user_tables
WHERE n_live_tup > 100
ORDER BY seq_scan DESC;

-- ============================================
-- 9. SUGGESTED INDEXES (Based on Usage)
-- ============================================

-- Find columns that might benefit from indexes
-- This is a simplified version - consider using pg_qualstats extension for better results
SELECT 
    'Consider indexing foreign keys without indexes' AS suggestion,
    conrelid::regclass AS table_name,
    a.attname AS column_name
FROM pg_constraint c
JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
WHERE c.contype = 'f'  -- Foreign key constraints
AND NOT EXISTS (
    SELECT 1
    FROM pg_index i
    WHERE i.indrelid = c.conrelid 
    AND a.attnum = ANY(i.indkey)
);

-- ============================================
-- 10. QUERY EXECUTION TIME DISTRIBUTION
-- ============================================

-- Active queries and their duration
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    CASE 
        WHEN state = 'active' THEN round(extract(epoch FROM (now() - query_start))::numeric, 2)
        ELSE NULL
    END AS query_duration_seconds,
    substring(query, 1, 100) AS current_query
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY query_start;

-- ============================================
-- PERFORMANCE RECOMMENDATIONS SUMMARY
-- ============================================

SELECT 'PERFORMANCE CHECK COMPLETE' AS status,
       'Review the results above for:' AS action_items,
       '1. Unused indexes to drop' AS item_1,
       '2. Tables needing VACUUM' AS item_2,
       '3. Missing indexes on foreign keys' AS item_3,
       '4. High sequential scan tables' AS item_4,
       '5. Cache hit rates below 99%' AS item_5;