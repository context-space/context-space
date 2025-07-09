# HashiCorp Vault Production Deployment

This directory contains the configuration and setup scripts for deploying HashiCorp Vault in a production environment for the Context Space application.

## Setup

The Vault setup includes:
1. A Docker Compose configuration for running Vault server
2. An initialization script for setting up Vault with the required configuration
3. The necessary policies for the Context Space application

## Directory Structure

```
hashivault/
├── config/            # Vault server configuration
│   └── config.hcl     # Main configuration file
├── policies/          # Predefined policy files
├── docker-compose.yml # Docker Compose configuration
├── init.sh            # Initialization script
└── README.md          # This file
```

## Deployment Instructions

### Prerequisites

- Docker and Docker Compose installed
- Network accessibility to the Vault instance from your application

### Step 1: Setup

Clone this repository and navigate to the hashivault directory:

```bash
git clone <repository-url>
cd <repository-path>/hashivault
```

### Step 2: Initialize Vault

Run the initialization script:

```bash
./init.sh
```

This script will:
1. Create the necessary directories
2. Start the Vault container
3. Initialize Vault and save the root token and unseal keys
4. Configure the transit engines for each region and credential type
5. Apply the application policy
6. Generate an application token

**IMPORTANT**: The script saves the root token and unseal keys to `vault_root_token.json`. In a production environment, you should secure this file or copy the keys to a secure location and delete the file.

## Accessing the Vault UI

The Vault UI is available at `http://localhost:8200/ui`. You can log in using the root token saved in `vault_root_token.json`.

## Operations

### Starting Vault

```bash
docker compose up -d
```

### Stopping Vault

```bash
docker compose down
```

### Checking Vault Status

```bash
docker compose exec vault vault status
```

### Creating a New Policy

1. Create a policy file in the `policies` directory
2. Apply the policy:

```bash
docker compose exec -e VAULT_TOKEN="<root-token>" vault vault policy write <policy-name> /vault/policies/<policy-file>
```

### Creating a New Token

```bash
docker compose exec -e VAULT_TOKEN="<root-token>" vault vault token create -policy=<policy-name>
```

## Backup and Restore

### Backup

Vault data is stored in the `data` directory. To backup Vault, stop the container and backup this directory:

```bash
docker compose down
tar -czf vault-backup-$(date +%Y%m%d).tar.gz data
docker compose up -d
```

### Restore

To restore from a backup:

```bash
docker compose down
rm -rf data
tar -xzf vault-backup-<date>.tar.gz
docker compose up -d
```

## Troubleshooting

### Connection Issues

If your application cannot connect to Vault, check:

1. Network connectivity between your application and Vault
2. Vault's status (it should be initialized and unsealed)
3. The token used by your application (it should have the correct policies)

### Vault is Sealed

If Vault becomes sealed (e.g., after a container restart), you need to unseal it:

```bash
# Get the unseal keys from vault_root_token.json
docker compose exec vault vault operator unseal <unseal-key-1>
docker compose exec vault vault operator unseal <unseal-key-2>
docker compose exec vault vault operator unseal <unseal-key-3>
```

### Authentication Issues

If you get authentication errors, check that you're using the correct token and that the token has the necessary policies.
