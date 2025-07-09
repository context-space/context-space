-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    sup_id UUID NOT NULL,
    email VARCHAR(255),
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_sup_id ON users(sup_id)
WHERE
    deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)
WHERE
    deleted_at IS NULL
    AND email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create user_infos table
CREATE TABLE IF NOT EXISTS user_infos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    info_metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_user_infos_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_user_infos_user_id ON user_infos(user_id);

CREATE INDEX IF NOT EXISTS idx_user_infos_deleted_at ON user_infos(deleted_at);

-- Create user_api_keys table
CREATE TABLE IF NOT EXISTS user_api_keys (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    key_value VARCHAR(64) NOT NULL,
    name VARCHAR(100),
    description TEXT,
    last_used TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_user_api_keys_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Add indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_api_keys_key_value ON user_api_keys(key_value)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_api_keys_user_id ON user_api_keys(user_id);

CREATE INDEX IF NOT EXISTS idx_user_api_keys_deleted_at ON user_api_keys(deleted_at);