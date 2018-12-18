package vaultlib

import (
	"encoding/json"
	"testing"
)

func Test_request_setJSONBody(t *testing.T) {
	var req request
	req.Method = "GET"
	payload := AppRoleCredentials{
		RoleID:   "a",
		SecretID: "b",
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		request request
		val     interface{}
		wantErr bool
	}{
		{"test", req, payload, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.setJSONBody(tt.val)
			if (err != nil) == tt.wantErr {
				var vaultAuth AppRoleCredentials
				_ = json.Unmarshal(tt.request.Body, &vaultAuth)
				if vaultAuth.RoleID != payload.RoleID {
					t.Errorf("not expected value")
				}
			}

		})
	}
}
