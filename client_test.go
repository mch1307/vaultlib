package vaultlib

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	appRoleCred := new(AppRoleCredentials)
	appRoleCred.RoleID = "abcd"
	appRoleCred.SecretID = "my-secret"
	tests := []struct {
		name string
		want Config
	}{
		{"DefaultConfig", Config{Address: "http://localhost:8200", InsecureSSL: true, Timeout: 30000000000, AppRoleCredentials: appRoleCred}},
		{"Custom", Config{Address: "http://localhost:8200", InsecureSSL: false, Timeout: 40000000000, CAPath: "/tmp", Token: "my-dev-root-vault-token", AppRoleCredentials: appRoleCred}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("VAULT_ROLEID", appRoleCred.RoleID)
			os.Setenv("VAULT_SECRETID", appRoleCred.SecretID)
			if tt.name == "Custom" {
				os.Setenv("VAULT_ADDR", "http://localhost:8200")
				os.Setenv("VAULT_SKIP_VERIFY", "0")
				os.Setenv("VAULT_CAPATH", "/tmp")
				os.Setenv("VAULT_TOKEN", "my-dev-root-vault-token")
				os.Setenv("VAULT_CLIENT_TIMEOUT", "40")

			}
			if got := NewConfig(); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, &tt.want)
			}
		})
	}
}

func TestConfig_SetAppRole(t *testing.T) {
	type fields struct {
		Address            string
		MaxRetries         int
		Timeout            time.Duration
		CAPath             string
		InsecureSSL        bool
		AppRoleCredentials AppRoleCredentials
		Token              string
	}
	var f fields
	type args struct {
		cred AppRoleCredentials
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"test1", f, args{AppRoleCredentials{RoleID: "role", SecretID: "secret"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Address:            tt.fields.Address,
				MaxRetries:         tt.fields.MaxRetries,
				Timeout:            tt.fields.Timeout,
				CAPath:             tt.fields.CAPath,
				InsecureSSL:        tt.fields.InsecureSSL,
				AppRoleCredentials: &tt.fields.AppRoleCredentials,
				Token:              tt.fields.Token,
			}
			if err := c.setAppRole(tt.args.cred); (err != nil) != tt.wantErr {
				t.Errorf("Config.setAppRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	// create client without token
	defaultCfg := NewConfig()
	vc, _ := NewClient(defaultCfg)
	// add token to client
	vc.Token = "my-dev-root-vault-token"
	// create new config with a vault token
	os.Setenv("VAULT_TOKEN", "my-dev-root-vault-token")
	cfg := NewConfig()
	os.Unsetenv("VAULT_TOKEN")
	// create new config without vault token
	wrongTokenConfig := NewConfig()
	wrongTokenConfig.Token = ""
	wrongTokenConfig.AppRoleCredentials.SecretID = "bad-secret"

	type args struct {
		c *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{"testOK", args{cfg}, vc, false},
		{"testFail", args{cfg}, vc, true},
		{"testNilConfig", args{wrongTokenConfig}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "testFail" {
				tt.args.c.Address = "hts://@\\x##ample.org:8080##@@"
			}
			got, err := NewClient(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !(got.Token == tt.want.Token) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewConfig() {
	myConfig := NewConfig()
	myConfig.Address = "http://localhost:8200"
}

func ExampleNewClient() {
	myConfig := NewConfig()
	myVaultClient, err := NewClient(myConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myVaultClient.Address)
}

func Example() {
	// Create a new config. Reads env variables, fallback to default value if needed
	vcConf := NewConfig()

	// Create new client
	vaultCli, err := NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}

	// Get the Vault KV secret from kv_v1/path/my-secret
	resV1, err := vaultCli.GetSecret("kv_v1/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resV1.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	// Get the Vault KVv2 secret kv_v2/path/my-secret
	resV2, err := vaultCli.GetSecret("kv_v2/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resV2.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	resJSON, err := vaultCli.GetSecret("kv_v2/path/json-secret")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("%v", resJSON.JSONSecret))
}
