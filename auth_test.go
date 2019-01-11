package vaultlib

import (
	"reflect"
	"testing"
	"time"
)

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
			c.setAppRole(tt.args.cred)
			if !reflect.DeepEqual(c.AppRoleCredentials, &tt.args.cred) {
				t.Errorf("Config.setAppRole() got %v, want %v", c.AppRoleCredentials, &tt.args.cred)
			}
		})
	}
}

// func TestVaultClient_setTokenFromAppRole(t *testing.T) {
// 	rightURL, _ := url.Parse("http://localhost:8200")
// 	badURL, _ := url.Parse("https://localhost:8200")
// 	conf := NewConfig()
// 	htCli := new(http.Client)
// 	type fields struct {
// 		Address    *url.URL
// 		HTTPClient *http.Client
// 		Config     *Config
// 		Token      string
// 		Status     string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{"tokenKO", fields{rightURL, htCli, conf, "bad-token", ""}, true},
// 		{"badUrl", fields{badURL, htCli, conf, "bad-token", ""}, true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				Address:    tt.fields.Address,
// 				HTTPClient: tt.fields.HTTPClient,
// 				Config:     tt.fields.Config,
// 				Token:      &vaultTokenInfo{ID: tt.fields.Token},
// 				Status:     tt.fields.Status,
// 			}
// 			if err := c.setTokenFromAppRole(); (err != nil) != tt.wantErr {
// 				t.Errorf("Client.setTokenFromAppRole() error = %v, wantErr %v", c.Token.ID, tt.fields.Token)
// 				//err, tt.wantErr)
// 			}
// 		})
// 	}
// }
