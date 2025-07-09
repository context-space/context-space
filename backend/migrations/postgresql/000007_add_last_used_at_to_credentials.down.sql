-- Drop index for last_used_at column
DROP INDEX IF EXISTS idx_credentials_last_used_at;

-- Drop last_used_at column from credentials table
ALTER TABLE credentials DROP COLUMN IF EXISTS last_used_at;