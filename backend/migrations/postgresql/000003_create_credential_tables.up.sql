-- Create credential type enum
CREATE TYPE credential_type AS ENUM ('oauth', 'apikey', 'basicauth');

-- Create credentials table
CREATE TABLE IF NOT EXISTS credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    provider_identifier VARCHAR(50) NOT NULL,
    credential_type credential_type NOT NULL,
    is_valid BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_credentials_user_id ON credentials(user_id);

CREATE INDEX IF NOT EXISTS idx_credentials_provider_identifier ON credentials(provider_identifier);

CREATE INDEX IF NOT EXISTS idx_credentials_credential_type ON credentials(credential_type);

CREATE INDEX IF NOT EXISTS idx_credentials_deleted_at ON credentials(deleted_at);

-- Create oauth_credentials table
CREATE TABLE IF NOT EXISTS oauth_credentials (
    credential_id UUID PRIMARY KEY,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    json_attributes JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_oauth_credentials_credential FOREIGN KEY (credential_id) REFERENCES credentials(id)
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_oauth_credentials_expiry ON oauth_credentials(expiry);

CREATE INDEX IF NOT EXISTS idx_oauth_credentials_deleted_at ON oauth_credentials(deleted_at);

-- Create apikey_credentials table
CREATE TABLE IF NOT EXISTS apikey_credentials (
    credential_id UUID PRIMARY KEY,
    json_attributes JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_apikey_credentials_credential FOREIGN KEY (credential_id) REFERENCES credentials(id)
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_apikey_credentials_deleted_at ON apikey_credentials(deleted_at);

-- Create oauth_states table
CREATE TABLE IF NOT EXISTS oauth_states (
    id UUID PRIMARY KEY,
    state VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    user_id UUID NOT NULL,
    provider_identifier VARCHAR(50) NOT NULL,
    json_attributes JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_oauth_states_state ON oauth_states(state);

CREATE INDEX IF NOT EXISTS idx_oauth_states_user_id ON oauth_states(user_id);

CREATE INDEX IF NOT EXISTS idx_oauth_states_deleted_at ON oauth_states(deleted_at);