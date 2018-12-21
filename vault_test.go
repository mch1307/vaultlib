package vaultlib

import (
	"fmt"
	"os"
	"testing"
)

func TestVaultClient_getKVVersion(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		os.Setenv("VAULT_ADDR", "http://localhost:8200")
		// Create a new config. Reads env variables, fallback to default value if needed
		conf := NewConfig()
		cred := AppRoleCredentials{
			RoleID:   "bb07197d-437f-9828-6512-94a5ec6c45a8",
			SecretID: "4d29046b-5b24-306b-e112-ec719fe5cd95",
		}
		_ = conf.SetAppRole(cred)
		cli, _ := NewClient(conf)

		err := cli.SetTokenFromAppRole()
		if err != nil {
			t.Errorf("error with app role auth: %v", err)
		}
		//cli.Token = "goodToken"

		gotVersion, name, err := cli.getKVInfo("kvv1/")
		fmt.Println("a", name)
		if err != nil {
			t.Errorf("Err:  %v", err)
		}
		if gotVersion != "1" {
			t.Errorf("VaultClient.getKVVersion() = %v, want %v", gotVersion, "1")
		}
	})

}
