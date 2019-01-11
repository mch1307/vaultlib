#!/bin/bash
#killall vault
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=my-dev-root-vault-token

nohup vault server -dev -dev-root-token-id ${VAULT_TOKEN}  > /dev/null 2>&1 &
# wait 5 seconds for Vault to be ready
sleep 5
# create KVs
vault secrets enable -path=kv_v1/path/ kv > /dev/null 2>&1
vault secrets enable -path=kv_v2/path/ kv > /dev/null 2>&1
vault kv enable-versioning kv_v2/path/ > /dev/null 2>&1

# create secrets
vault kv put kv_v1/path/my-secret my-v1-secret=my-v1-secret-value > /dev/null 2>&1
vault kv put kv_v2/path/my-secret my-first-secret=my-first-secret-value my-second-secret=my-second-secret-value > /dev/null 2>&1

# create policy
vault policy write VaultDevAdmin ./VaultPolicy.hcl > /dev/null 2>&1

# create approle
vault auth enable approle > /dev/null 2>&1
vault write auth/approle/role/my-role policies=VaultDevAdmin secret_id_ttl=1000m token_num_uses=5 token_ttl=10s token_max_ttl=30m secret_id_num_uses=40 > /dev/null 2>&1
export VAULT_ROLEID=`vault read -field=role_id auth/approle/role/my-role/role-id`
export VAULT_SECRETID=`vault write -field=secret_id -f auth/approle/role/my-role/secret-id`

unset VAULT_TOKEN
