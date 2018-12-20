package vaultlib

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	appRoleCred := new(AppRoleCredentials)
	appRoleCred.RoleID = "abcd"
	tests := []struct {
		name string
		want Config
	}{
		{"DefaultConfig", Config{Address: "http://localhost:8200", InsecureSSL: true, Timeout: 30000000000, AppRoleCredentials: appRoleCred}},
		{"Custom", Config{Address: "https://localhost:8200", InsecureSSL: false, Timeout: 40000000000, CAPath: "/tmp", Token: "abcd", AppRoleCredentials: appRoleCred}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("VAULT_ROLEID", "abcd")
			if tt.name == "Custom" {
				os.Setenv("VAULT_ADDR", "https://localhost:8200")
				os.Setenv("VAULT_SKIP_VERIFY", "0")
				os.Setenv("VAULT_CAPATH", "/tmp")
				os.Setenv("VAULT_TOKEN", "abcd")
				os.Setenv("VAULT_ROLE_ID", "abcd")
				os.Setenv("VAULT_CLIENT_TIMEOUT", "40")

			}
			if got := NewConfig(); !reflect.DeepEqual(got, &tt.want) {
				fmt.Println(tt.want.AppRoleCredentials.RoleID)
				fmt.Println(got.AppRoleCredentials.RoleID)
				t.Errorf("NewConfig() = %v, want %v %v", got, &tt.want, got.AppRoleCredentials.RoleID)
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
			if err := c.SetAppRole(&tt.args.cred); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetAppRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	// create client without token
	defaultCfg := NewConfig()
	vc, _ := NewClient(defaultCfg)
	// add token to client
	vc.Token = "abcd"
	// create new config with a vault token
	os.Setenv("VAULT_TOKEN", "abcd")
	cfg := NewConfig()

	type args struct {
		c *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *VaultClient
		wantErr bool
	}{
		{"testOK", args{cfg}, vc, false},
		{"testFail", args{cfg}, vc, true},
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
