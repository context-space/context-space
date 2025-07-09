#!/bin/bash
set -e

echo "===== Vault Setup - Phase 1: Initialization ====="
echo "Creating necessary directories..."
mkdir -p /root/data/vault/data /root/data/vault/logs /root/data/vault/policies

# 修改目录权限，确保能被容器内的vault用户访问
echo "Setting directory permissions..."
chmod 777 /root/data/vault/data /root/data/vault/logs /root/data/vault/policies

# 确保策略文件就位
echo "Copying policy file..."
cp policies/app.hcl /root/data/vault/policies/

# 启动Vault
echo "Starting Vault..."
docker compose up -d

echo "Waiting for Vault to start..."
sleep 10

# 检查Vault状态
echo "Checking Vault status..."
INIT_STATUS=$(docker compose exec vault vault status -format=json 2>/dev/null || echo '{"initialized": false}')
INITIALIZED=$(echo $INIT_STATUS | grep -o '"initialized":[^,}]*' | cut -d: -f2 | tr -d ' "')

if [ "$INITIALIZED" != "true" ]; then
    echo "===== Vault Setup - Phase 2: Manual Initialization Required ====="
    echo "Vault needs to be initialized. Please run the following commands:"
    echo ""
    echo "# Initialize Vault and save keys to the host machine"
    echo "docker compose exec vault vault operator init -format=json > vault_root_token.json"
    echo "cat vault_root_token.json  # Save these keys securely!"
    echo ""
    echo "# Unseal Vault with three keys from the output"
    echo "docker compose exec vault vault operator unseal  # Enter the first Unseal Key when prompted"
    echo "docker compose exec vault vault operator unseal  # Enter the second Unseal Key when prompted"
    echo "docker compose exec vault vault operator unseal  # Enter the third Unseal Key when prompted"
    echo ""
    echo "# Verify that Vault is unsealed"
    echo "docker compose exec vault vault status"
    echo ""
    echo "===== Vault Setup - Phase 3: Configure Policies and Transit Engines ====="
    echo "After initializing and unsealing, configure Vault with the following:"
    echo ""
    echo "# Copy the setup script to container"
    echo "docker compose cp setup-policies.sh vault:/tmp/"
    echo "docker compose exec vault chmod +x /tmp/setup-policies.sh"
    echo ""
    echo "# Run the setup script with your root token"
    echo "docker compose exec -e VAULT_TOKEN=\$ROOT_TOKEN vault /tmp/setup-policies.sh"
else
    echo "===== Vault Setup - Phase 2: Already Initialized ====="
    echo "Vault is already initialized. If sealed, you'll need to unseal it with:"
    echo ""
    echo "docker compose exec vault vault operator unseal  # Enter the first Unseal Key when prompted"
    echo "docker compose exec vault vault operator unseal  # Enter the second Unseal Key when prompted"
    echo "docker compose exec vault vault operator unseal  # Enter the third Unseal Key when prompted"
    echo ""
    echo "===== Vault Setup - Phase 3: Configure Policies and Transit Engines ====="
    echo "To configure or update policies and transit engines, run:"
    echo ""
    echo "# With your root token, run the setup script"
    echo "docker compose cp setup-policies.sh vault:/tmp/"
    echo "docker compose exec vault chmod +x /tmp/setup-policies.sh"
    echo ""
    echo "# Run the setup script with your root token"
    echo "docker compose exec -e VAULT_TOKEN=<your-root-token> vault /tmp/setup-policies.sh"
fi

echo ""
echo "===== IMPORTANT NOTES ====="
echo "1. Store your Unseal Keys and Root Token securely - they cannot be recovered if lost!"
echo "2. In production, do not store these keys in plaintext or in version control"
echo "3. Your app should use the app token, not the root token"
echo "4. The setup-policies.sh script can be run multiple times safely"
echo "5. All tokens and keys are now saved on your host machine in vault_root_token.json and app_token.json"
echo ""
echo "For more information, see the README.md file"
