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
				Token:      &vaultTokenInfo{ID: tt.fields.Token},
				Status:     tt.fields.Status,
			}
			if err := c.setTokenFromAppRole(); (err != nil) != tt.wantErr {
				t.Errorf("Client.setTokenFromAppRole() error = %v, wantErr %v", c.Token.ID, tt.fields.Token)
				//err, tt.wantErr)
			}
		})
	}
}
