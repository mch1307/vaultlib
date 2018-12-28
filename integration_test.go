package vaultlib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

var vaultRoleID, vaultSecretID string

func prepareVault() {
	go startVault()
	// wait 30 seconds vault
	time.Sleep(30 * time.Second)
	cmd := exec.Command("./vault", "read", "-field=role_id", "auth/approle/role/my-role/role-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("error getting role id %v %v", err, out)
	}
	vaultRoleID = string(out)

	cmd = exec.Command("./vault", "write", "-field=secret_id", "-f", "auth/approle/role/my-role/secret-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")
	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("error getting secret id %v", err)
	}
	vaultSecretID = string(out)
	os.Unsetenv("VAULT_TOKEN")

}

func startVault() {
	cmd := exec.Command("./initVaultDev.sh")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Wait()

}
func TestMain(m *testing.M) {
	fmt.Println("TestMain: Preparing Vault server")
	prepareVault()
	ret := m.Run()
	os.Exit(ret)
}

func TestVaultClient_GetVaultSecret(t *testing.T) {

	conf := NewConfig()
	conf.AppRoleCredentials.RoleID = vaultRoleID
	conf.AppRoleCredentials.SecretID = vaultSecretID
	vc, err := NewClient(conf)
	if err != nil {
		t.Errorf("Failed to get vault cli %v", err)
	}

	tests := []struct {
		name    string
		cli     *VaultClient
		path    string
		wantKv  map[string]string
		wantErr bool
	}{
		{"kvv1", vc, "kv_v1/path/my-secret", map[string]string{"my-v1-secret": "my-v1-secret-value"}, false},
		{"kvv2", vc, "kv_v2/path/my-secret", map[string]string{"my-first-secret": "my-first-secret-value",
			"my-second-secret": "my-second-secret-value"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cli
			gotKv, err := c.GetVaultSecret(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("VaultClient.GetVaultSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKv, tt.wantKv) {
				t.Errorf("VaultClient.GetVaultSecret() = %v, want %v", gotKv, tt.wantKv)
			}
		})
	}
}
