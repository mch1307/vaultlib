package vaultlib

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func Test_request_setJSONBody(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		var cred AppRoleCredentials
		cred.RoleID = "aa"
		cred.SecretID = "bb"
		htCli := new(http.Client)
		url, _ := url.Parse("http://localhot:8200")
		req, _ := newRequest("GET", "", url, htCli)
		err := req.setJSONBody(cred)

		var vaultAuth AppRoleCredentials
		decoder := json.NewDecoder(req.Req.Body)
		err = decoder.Decode(&vaultAuth)
		if err != nil {
			t.Error("error parsing")
		}
		if vaultAuth.RoleID != "aa" {
			t.Errorf("got %v expecting aa", vaultAuth.RoleID)
		}

	})

}
