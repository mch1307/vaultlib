#!/bin/bash
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=my-dev-root-vault-token
export VAULT_VERSION=${1:-1.0.1}

CURRENT_VAULT=`./vault version | cut -d'v' -f2 | cut -d' ' -f1`
if [ "$CURRENT_VAULT" != "$VAULT_VERSION" ]; then
    rm -rf ./vault
    curl -kO https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip
    unzip vault_${VAULT_VERSION}_linux_amd64.zip
fi

./vault server -dev -dev-root-token-id ${VAULT_TOKEN}  > /tmp/vaultdev.log &
# wait for vault server to be ready
sleep 5

# create KVs
./vault secrets enable -path=kv_v1/path/ kv >> /tmp/vaultdev.log
./vault secrets enable -path=kv_v2/path/ kv >> /tmp/vaultdev.log
./vault kv enable-versioning kv_v2/path/ >> /tmp/vaultdev.log

# create secrets
./vault kv put kv_v1/path/my-secret my-v1-secret=my-v1-secret-value >> /tmp/vaultdev.log
./vault kv put kv_v2/path/my-secret my-first-secret=my-first-secret-value my-second-secret=my-second-secret-value >> /tmp/vaultdev.log
./vault kv put kv_v2/path/json-secret @./test-files/secret.json >> /tmp/vaultdev.log
./vault kv put kv_v1/path/json-secret @./test-files/secret.json >> /tmp/vaultdev.log

# create policy
./vault policy write VaultDevAdmin test-files/VaultPolicy.hcl >> /tmp/vaultdev.log

# create approle
./vault auth enable approle >> /tmp/vaultdev.log
#./vault write auth/approle/role/my-role policies=VaultDevAdmin secret_id_ttl=100m token_num_uses=100 token_ttl=100m token_max_ttl=300m secret_id_num_uses=40 >> /tmp/vaultdev.log
./vault write auth/approle/role/my-role policies=VaultDevAdmin token_num_uses=100 token_ttl=360s token_max_ttl=300m secret_id_num_uses=40 >> /tmp/vaultdev.log

unset VAULT_TOKEN
