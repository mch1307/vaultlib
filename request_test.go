package vaultlib

import (
	"encoding/json"
	"net/url"
	"testing"
)

func Test_request_setJSONBody(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		var cred AppRoleCredentials
		cred.RoleID = "aa"
		cred.SecretID = "bb"
		var req request
		req.URL, _ = url.Parse("http://localhot:8200")
		req.prepareRequest()
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
