package vaultlib

import (
	"testing"
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
