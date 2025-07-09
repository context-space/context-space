# Policy for the application to access transit encryption keys
# AUTOMATICALLY GENERATED - DO NOT EDIT MANUALLY
# Generated on Tue May 13 18:17:00 HKT 2025

# Region-specific transit paths for oauth credentials in eu
path "transit-eu-oauth/encrypt/oauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-oauth/decrypt/oauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-oauth/rewrap/oauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-oauth/keys/oauth-creds-eu-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for apikey credentials in eu
path "transit-eu-apikey/encrypt/apikey-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-apikey/decrypt/apikey-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-apikey/rewrap/apikey-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-apikey/keys/apikey-creds-eu-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for basicauth credentials in eu
path "transit-eu-basicauth/encrypt/basicauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-basicauth/decrypt/basicauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-basicauth/rewrap/basicauth-creds-eu-key" {
  capabilities = ["update"]
}

path "transit-eu-basicauth/keys/basicauth-creds-eu-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for oauth credentials in us
path "transit-us-oauth/encrypt/oauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-oauth/decrypt/oauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-oauth/rewrap/oauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-oauth/keys/oauth-creds-us-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for apikey credentials in us
path "transit-us-apikey/encrypt/apikey-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-apikey/decrypt/apikey-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-apikey/rewrap/apikey-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-apikey/keys/apikey-creds-us-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for basicauth credentials in us
path "transit-us-basicauth/encrypt/basicauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-basicauth/decrypt/basicauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-basicauth/rewrap/basicauth-creds-us-key" {
  capabilities = ["update"]
}

path "transit-us-basicauth/keys/basicauth-creds-us-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for oauth credentials in cn
path "transit-cn-oauth/encrypt/oauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-oauth/decrypt/oauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-oauth/rewrap/oauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-oauth/keys/oauth-creds-cn-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for apikey credentials in cn
path "transit-cn-apikey/encrypt/apikey-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-apikey/decrypt/apikey-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-apikey/rewrap/apikey-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-apikey/keys/apikey-creds-cn-key" {
  capabilities = ["read"]
}

# Region-specific transit paths for basicauth credentials in cn
path "transit-cn-basicauth/encrypt/basicauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-basicauth/decrypt/basicauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-basicauth/rewrap/basicauth-creds-cn-key" {
  capabilities = ["update"]
}

path "transit-cn-basicauth/keys/basicauth-creds-cn-key" {
  capabilities = ["read"]
}

# Admin users can rotate keys for all transit paths
path "transit-*/keys/*/rotate" {
  capabilities = ["update"]
}
