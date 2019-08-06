package vaultlib

import (
	"net/http"
	"net/url"
	"testing"
)

func TestVaultClient_setTokenFromAppRole(t *testing.T) {
	rightURL, _ := url.Parse("http://localhost:8200")
	badURL, _ := url.Parse("https://localhost:8200")
	conf := NewConfig()
	anyMountPoint := "anyMountPoint"
	anyCreds := NewConfig().AppRoleCredentials
	anyCreds.MountPoint = anyMountPoint
	htCli := new(http.Client)
	type fields struct {
		Address            *url.URL
		HTTPClient         *http.Client
		AppRoleCredentials *AppRoleCredentials
		//Config     *Config
		Token  string
		Status string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"tokenKO",
			fields{
				rightURL,
				htCli,
				conf.AppRoleCredentials,
				"bad-token",
				""},
			true},
		{"badUrl",
			fields{
				badURL,
				htCli,
				conf.AppRoleCredentials,
				"bad-token",
				""},
			true},
		{"badMountPoint",
			fields{
				rightURL,
				htCli,
				anyCreds,
				"bad-token",
				""},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				address:            tt.fields.Address,
				httpClient:         tt.fields.HTTPClient,
				appRoleCredentials: tt.fields.AppRoleCredentials,
				//config:     tt.fields.Config,
				token:  &VaultTokenInfo{ID: tt.fields.Token},
				status: tt.fields.Status,
			}
			if err := c.setTokenFromAppRole(); (err != nil) != tt.wantErr {
				t.Errorf("Client.setTokenFromAppRole() error = %v, wantErr %v", c.token.ID, tt.fields.Token)
				//err, tt.wantErr)
			}
		})
	}
}
