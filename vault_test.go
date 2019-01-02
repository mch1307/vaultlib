package vaultlib

import (
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestVaultClient_getKVInfo(t *testing.T) {
	conf := NewConfig()
	conf.Address = "http://localhost:8200"
	cred := AppRoleCredentials{
		RoleID:   vaultRoleID,
		SecretID: vaultSecretID,
	}
	_ = conf.setAppRole(cred)
	badReqConf := NewConfig()
	badReqConf.Address = "https://localhost:8200"
	noCred := AppRoleCredentials{
		RoleID:   "",
		SecretID: "",
	}
	noCredConf := NewConfig()
	noCredConf.AppRoleCredentials = &noCred

	type fields struct {
		Config *Config
	}
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantVersion string
		wantName    string
		wantErr     bool
	}{
		{"foundV1", fields{conf}, args{"kv_v1/path/my-secret"}, "1", "kv_v1/path/", false},
		{"foundV2", fields{conf}, args{"kv_v2/path/my-secret"}, "2", "kv_v2/path/", false},
		{"notFound", fields{conf}, args{"notExist/my-secret"}, "", "", true},
		{"badRequest", fields{badReqConf}, args{"notExist/my-secret"}, "", "", true},
		{"NoCred", fields{noCredConf}, args{"notExist/my-secret"}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewClient(tt.fields.Config)
			gotVersion, gotName, err := c.getKVInfo(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("VaultClient.getKVInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("VaultClient.getKVInfo() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
			if gotName != tt.wantName {
				t.Errorf("VaultClient.getKVInfo() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestVaultClient_setTokenFromAppRole(t *testing.T) {
	rightURL, _ := url.Parse("http://localhost:8200")
	badURL, _ := url.Parse("https://localhost:8200")
	conf := NewConfig()
	htCli := new(http.Client)
	type fields struct {
		Address    *url.URL
		HTTPClient *http.Client
		Config     *Config
		Token      string
		Status     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"tokenKO", fields{rightURL, htCli, conf, "bad-token", ""}, true},
		{"badUrl", fields{badURL, htCli, conf, "bad-token", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &VaultClient{
				Address:    tt.fields.Address,
				HTTPClient: tt.fields.HTTPClient,
				Config:     tt.fields.Config,
				Token:      tt.fields.Token,
				Status:     tt.fields.Status,
			}
			if err := c.setTokenFromAppRole(); (err != nil) != tt.wantErr {
				t.Errorf("VaultClient.setTokenFromAppRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVaultClient_GetVaultSecret(t *testing.T) {
	_ = os.Unsetenv("VAULT_TOKEN")
	conf := NewConfig()
	conf.AppRoleCredentials.RoleID = vaultRoleID
	conf.AppRoleCredentials.SecretID = vaultSecretID
	vc, err := NewClient(conf)
	if err != nil {
		t.Errorf("Failed to get vault cli %v", err)
	}
	conf.Address = "https://localhost:8200"
	badCli, _ := NewClient(conf)

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
		{"invalidURL", badCli, "kv_v1/path/my-secret", map[string]string{}, true},
	}
	//wait so that token renewal takes place
	time.Sleep(12 * time.Second)
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
