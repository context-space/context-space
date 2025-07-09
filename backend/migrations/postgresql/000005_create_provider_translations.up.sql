-- Create provider_translations table
CREATE TABLE IF NOT EXISTS provider_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_identifier VARCHAR(255) NOT NULL,
    language_code VARCHAR(10) NOT NULL,
    translations JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_provider_translations_identifier_language_code ON provider_translations(provider_identifier, language_code)