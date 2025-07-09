-- Create invocations table
CREATE TABLE IF NOT EXISTS invocations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    provider_identifier VARCHAR(50) NOT NULL,
    operation_identifier VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    duration BIGINT NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    json_attributes JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_invocations_user_id ON invocations(user_id);

CREATE INDEX IF NOT EXISTS idx_invocations_provider_identifier ON invocations(provider_identifier);

CREATE INDEX IF NOT EXISTS idx_invocations_operation_identifier ON invocations(operation_identifier);

CREATE INDEX IF NOT EXISTS idx_invocations_deleted_at ON invocations(deleted_at);