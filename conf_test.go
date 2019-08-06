package vaultlib

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	appRoleCred := new(AppRoleCredentials)
	appRoleCred.RoleID = "abcd"
	appRoleCred.SecretID = "my-secret"
	appRoleCred.MountPoint = "approle"
	tests := []struct {
		name string
		want Config
	}{
		{"DefaultConfig", Config{Address: "http://localhost:8200", InsecureSSL: true, Timeout: 30000000000, AppRoleCredentials: appRoleCred}},
		{"Custom", Config{Address: "http://localhost:8200", InsecureSSL: false, Timeout: 40000000000, CACert: "/tmp", Token: "my-dev-root-vault-token", AppRoleCredentials: appRoleCred}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("VAULT_ROLEID", appRoleCred.RoleID)
			os.Setenv("VAULT_SECRETID", appRoleCred.SecretID)
			os.Setenv("VAULT_MOUNTPOINT", appRoleCred.MountPoint)
			if tt.name == "Custom" {
				os.Setenv("VAULT_ADDR", "http://localhost:8200")
				os.Setenv("VAULT_SKIP_VERIFY", "0")
				os.Setenv("VAULT_CACERT", "/tmp")
				os.Setenv("VAULT_TOKEN", "my-dev-root-vault-token")
				os.Setenv("VAULT_CLIENT_TIMEOUT", "40")
			}
			if got := NewConfig(); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, &tt.want)
			}
		})
	}
	//clean env
	os.Unsetenv("VAULT_ROLEID")
	os.Unsetenv("VAULT_SECRETID")
	os.Unsetenv("VAULT_MOUNTPOINT")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_SKIP_VERIFY")
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULT_CLIENT_TOKEN")
}

func ExampleNewConfig() {
	myConfig := NewConfig()
	myConfig.Address = "http://localhost:8200"
}
