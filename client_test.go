package vaultlib

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	// create client without token
	defaultCfg := NewConfig()
	defaultCfg.AppRoleCredentials.RoleID = vaultRoleID
	defaultCfg.AppRoleCredentials.SecretID = vaultSecretID
	vc, _ := NewClient(defaultCfg)
	// create new config with a vault token
	os.Setenv("VAULT_TOKEN", "my-renewable-token")
	cfg := NewConfig()
	// create new config without vault token
	os.Unsetenv("VAULT_TOKEN")
	wrongTokenConfig := NewConfig()
	wrongTokenConfig.Token = ""
	wrongTokenConfig.AppRoleCredentials.SecretID = "bad-secret"
	wrongTokenConfig.AppRoleCredentials.RoleID = "bad-roleid"
	noAppRoleConfig := NewConfig()
	noAppRoleConfig.AppRoleCredentials.RoleID = ""
	noAppRoleConfig.Token = "bad-token"

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
		{"testNoCfg", args{}, nil, true},
		{"testFail", args{cfg}, vc, true},
		{"testNilConfig", args{wrongTokenConfig}, nil, true},
		{"noAppRoleConfig", args{noAppRoleConfig}, nil, true},
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
			if !tt.wantErr && !(got.status == tt.want.status) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Example() {
	// Create a new config. Reads env variables, fallback to default value if needed
	vcConf := NewConfig()

	// Add the Vault approle secretid after having read it from docker secret
	// vcConf.AppRoleCredentials.SecretID

	// Create new client
	vaultCli, err := NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}

	// Get the Vault KV secret from kv_v1/path/my-secret
	kvV1, err := vaultCli.GetSecret("kv_v1/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range kvV1.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	// Get the Vault KVv2 secret kv_v2/path/my-secret
	kvV2, err := vaultCli.GetSecret("kv_v2/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range kvV2.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	jsonSecret, err := vaultCli.GetSecret("kv_v2/path/json-secret")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("%v", jsonSecret.JSONSecret))
}

func ExampleNewClient() {
	myConfig := NewConfig()
	myVaultClient, err := NewClient(myConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myVaultClient.address)
}

func ExampleClient_IsAuthenticated() {
	myConfig := NewConfig()
	myVaultClient, err := NewClient(myConfig)
	if err != nil {
		fmt.Println(err)
	}
	if myVaultClient.IsAuthenticated() {
		fmt.Println("myVaultClient's connection is ok")
	}
}

func TestClient_IsAuthenticated(t *testing.T) {
	conf := NewConfig()
	conf.Token = "my-dev-root-vault-token"
	authCli, _ := NewClient(conf)
	conf.Token = "bad-token"
	badCli, _ := NewClient(conf)
	tests := []struct {
		name string
		cli  *Client
		want bool
	}{
		{"auth", authCli, true},
		{"noAuth", badCli, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cli
			if got := c.IsAuthenticated(); got != tt.want {
				t.Errorf("Client.IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetTokenInfo(t *testing.T) {
	defaultCfg := NewConfig()
	defaultCfg.Token = "my-dev-root-vault-token"
	client, _ := NewClient(defaultCfg)
	tokenOK := new(VaultTokenInfo)
	tokenOK.ID = defaultCfg.Token
	tests := []struct {
		name string
		cli  *Client
		want *VaultTokenInfo
	}{
		{"OK", client, tokenOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client
			if got := c.GetTokenInfo().ID; !reflect.DeepEqual(got, tt.want.ID) {
				t.Errorf("Client.GetTokenInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetStatus(t *testing.T) {
	defaultCfg := NewConfig()
	defaultCfg.Token = "my-dev-root-vault-token"
	client, _ := NewClient(defaultCfg)
	tests := []struct {
		name string
		cli  *Client
		want string
	}{
		{"Ready", client, "Token ready"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client
			if got := c.GetStatus(); got != tt.want {
				t.Errorf("Client.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
