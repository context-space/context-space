-- Drop tables in reverse order (dependencies first)
DROP TABLE IF EXISTS apikey_credentials;

DROP TABLE IF EXISTS oauth_credentials;

DROP TABLE IF EXISTS oauth_states;

DROP TABLE IF EXISTS credentials;

-- Drop enum type
DROP TYPE IF EXISTS credential_type;