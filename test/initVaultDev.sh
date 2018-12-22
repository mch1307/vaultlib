#!/bin/bash
killall vault
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=my-dev-root-vault-token

vault server -dev -dev-root-token-id ${VAULT_TOKEN} &

# create KVs
vault secrets enable -path=kv_v1/path/ kv
vault secrets enable -path=kv_v2/path/ kv
vault kv enable-versioning kv_v2/path/

# create secrets
vault kv put kv_v1/path/my-secret secret=value
vault kv put kv_v2/path/my-secret secret=value
vault kv put kv_v2/path/my-secret secret2=value2

# create policy
vault policy write VaultDevAdmin vaultDevAdminPolicy.hcl

# create approle
vault auth enable approle
vault write auth/approle/role/my-role policies=VaultDevAdmin secret_id_ttl=1000m token_num_uses=100 token_ttl=1000m token_max_ttl=3000m secret_id_num_uses=40
export VAULT_ROLEID=`vault read -field=role_id auth/approle/role/my-role/role-id`
export VAULT_SECRETID=`vault write -field=secret_id -f auth/approle/role/my-role/secret-id`

