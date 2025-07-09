#!/bin/sh
# 此脚本应在 Vault 容器内运行

set -e

# 使用root token
if [ -z "$VAULT_TOKEN" ]; then
    echo "Error: VAULT_TOKEN environment variable must be set"
    echo "Usage: VAULT_TOKEN=<your-root-token> ./setup-policies.sh"
    exit 1
fi

# 设置Vault地址
export VAULT_ADDR=http://127.0.0.1:8200

# 验证token有效性
echo "Verifying token..."
vault token lookup >/dev/null || {
    echo "Error: Invalid token or Vault is not accessible"
    exit 1
}

# 定义区域和凭证类型
REGIONS="eu us cn"
CRED_TYPES="oauth apikey basicauth"

# 启用transit引擎并创建密钥
for region in $REGIONS; do
    for cred_type in $CRED_TYPES; do
        # 创建transit挂载点
        MOUNT_PATH="transit-${region}-${cred_type}"

        echo "Enabling transit engine at $MOUNT_PATH"
        vault secrets enable -path=$MOUNT_PATH transit || echo "Mount $MOUNT_PATH already exists"

        # 为该挂载点创建密钥
        KEY_NAME="${cred_type}-creds-${region}-key"
        echo "Creating key $KEY_NAME at $MOUNT_PATH"
        vault write -f "$MOUNT_PATH/keys/$KEY_NAME" || echo "Key $KEY_NAME already exists"
    done
done

# 应用策略文件
if [ -f "/vault/policies/app.hcl" ]; then
    echo "Applying policy from /vault/policies/app.hcl..."
    vault policy write app /vault/policies/app.hcl

    # 创建应用令牌
    echo "Creating app token..."
    APP_TOKEN=$(vault token create -policy=app -display-name=context-space-app -format=json)
    echo "$APP_TOKEN" >/vault/policies/app_token.json
    echo "Application token saved to /vault/policies/app_token.json"

    # 展示应用令牌信息
    echo "App token client_token: $(echo $APP_TOKEN | grep -o '"client_token":"[^"]*"' | cut -d':' -f2 | tr -d '"' | tr -d ',')"
else
    echo "Error: Policy file not found at /vault/policies/app.hcl"
    exit 1
fi

echo "Policy and transit setup complete!"
