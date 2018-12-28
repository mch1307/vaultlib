package vaultlib

import (
	"os"
	"testing"
)

func TestVaultClient_getKVVersion(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		os.Setenv("VAULT_ADDR", "http://localhost:8200")
		// Create a new config. Reads env variables, fallback to default value if needed
		conf := NewConfig()
		cred := AppRoleCredentials{
			RoleID:   vaultRoleID,
			SecretID: vaultSecretID,
		}
		_ = conf.setAppRole(cred)
		cli, _ := NewClient(conf)

		err := cli.setTokenFromAppRole()
		if err != nil {
			t.Errorf("error with app role auth: %v", err)
		}

		gotVersion, _, err := cli.getKVInfo("kv_v1/path/")
		if err != nil {
			t.Errorf("Err:  %v", err)
		}
		if gotVersion != "1" {
			t.Errorf("VaultClient.getKVVersion() = %v, want %v", gotVersion, "1")
		}
	})

}
