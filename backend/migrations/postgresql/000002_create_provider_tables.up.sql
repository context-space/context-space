-- Create providers table
CREATE TABLE IF NOT EXISTS providers (
    id UUID PRIMARY KEY,
    identifier VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    auth_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    icon_url TEXT,
    json_attributes JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add unique index for identifier
CREATE UNIQUE INDEX IF NOT EXISTS idx_providers_identifier ON providers(identifier)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_providers_deleted_at ON providers(deleted_at);

-- Create operations table
CREATE TABLE IF NOT EXISTS operations (
    id UUID PRIMARY KEY,
    identifier VARCHAR(50) NOT NULL,
    provider_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    json_attributes JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_operations_provider FOREIGN KEY (provider_id) REFERENCES providers(id)
);

-- Add unique composite index for identifier and provider_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_operations_identifier_provider ON operations(identifier, provider_id)
WHERE
    deleted_at IS NULL;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_operations_provider_id ON operations(provider_id);

CREATE INDEX IF NOT EXISTS idx_operations_deleted_at ON operations(deleted_at);