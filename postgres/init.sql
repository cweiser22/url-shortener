CREATE TABLE IF NOT EXISTS url_access_log (
    short_code TEXT PRIMARY KEY,
    visit_count BIGINT DEFAULT 1,
    last_access TIMESTAMPTZ DEFAULT NOW()
);

-- Ensure an index on last_access for efficient time-based lookups
CREATE INDEX IF NOT EXISTS idx_last_access ON url_access_log (last_access);
