package vaultlib

import (
	"encoding/json"
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
				t.Errorf("Client.getKVInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("Client.getKVInfo() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
			if gotName != tt.wantName {
				t.Errorf("Client.getKVInfo() gotName = %v, want %v", gotName, tt.wantName)
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
			c := &Client{
				Address:    tt.fields.Address,
				HTTPClient: tt.fields.HTTPClient,
				Config:     tt.fields.Config,
				Token:      tt.fields.Token,
				Status:     tt.fields.Status,
			}
			if err := c.setTokenFromAppRole(); (err != nil) != tt.wantErr {
				t.Errorf("Client.setTokenFromAppRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVaultClient_GetSecret(t *testing.T) {
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
	expectedJSON := []byte(`{"json-secret":{"first-secret":"first-value","second-secret":"second-value"}}`)
	tests := []struct {
		name     string
		cli      *Client
		path     string
		wantKv   map[string]string
		wantJSON json.RawMessage
		wantErr  bool
	}{
		{"kvv1", vc, "kv_v1/path/my-secret", map[string]string{"my-v1-secret": "my-v1-secret-value"}, nil, false},
		{"kvv2", vc, "kv_v2/path/my-secret", map[string]string{"my-first-secret": "my-first-secret-value",
			"my-second-secret": "my-second-secret-value"}, nil, false},
		{"json-secretV2", vc, "kv_v2/path/json-secret", map[string]string{}, expectedJSON, false},
		{"json-secretV1", vc, "kv_v1/path/json-secret", map[string]string{}, expectedJSON, false},
		{"invalidURL", badCli, "kv_v1/path/my-secret", map[string]string{}, nil, true},
	}
	//wait so that token renewal takes place
	time.Sleep(12 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cli
			res, err := c.GetSecret(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res.KV, tt.wantKv) || !reflect.DeepEqual(res.JSONSecret, tt.wantJSON) {
				t.Errorf("Client.GetSecret() = %v, want %v", res.KV, tt.wantKv)
			}
		})
	}
}
