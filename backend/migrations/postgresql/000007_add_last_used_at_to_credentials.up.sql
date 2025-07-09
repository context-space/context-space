-- Add last_used_at column to credentials table
ALTER TABLE credentials ADD COLUMN last_used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add index for last_used_at column
CREATE INDEX IF NOT EXISTS idx_credentials_last_used_at ON credentials(last_used_at); 