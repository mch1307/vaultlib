package vaultlib

import (
	"os"
	"testing"
)

func TestVaultClient_getKVVersion(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		os.Setenv("VAULT_ADDR", "http://localhost:8500")
		// Create a new config. Reads env variables, fallback to default value if needed
		conf := NewConfig()
		cli, _ := NewClient(conf)
		cli.Token = "goodToken"

		gotVersion, err := cli.getKVVersion("kv1")
		if err != nil {
			t.Errorf("Err:  %v", err)
		}
		if gotVersion != "1" {
			t.Errorf("VaultClient.getKVVersion() = %v, want %v", gotVersion, "1")
		}
	})

}
