package vaultlib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func Test_request_setJSONBody(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		var cred AppRoleCredentials
		cred.RoleID = "aa"
		cred.SecretID = "bb"
		//htCli := new(http.Client)
		url, _ := url.Parse("http://localhot:8200")
		req, _ := newRequest("GET", "", url)
		err := req.setJSONBody(cred)
		if err != nil {
			fmt.Println(err)
		}

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

// func Test_newRequest(t *testing.T) {
// 	testUrl, _ := url.Parse("http://localhost:8200")
// 	type args struct {
// 		method string
// 		token  string
// 		url    *url.URL
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *request
// 		wantErr bool
// 	}{
// 		{"err", args{"\\@", "gcccvk", testUrl}, nil, true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := newRequest(tt.args.method, tt.args.token, tt.args.url)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("newRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_request_execute(t *testing.T) {
	nReq, _ := http.NewRequest("GET", "hppT:/localhost.18200", nil)
	htcli := cleanhttp.DefaultPooledClient()
	var rsp vaultResponse
	type fields struct {
		Req        *http.Request
		HTTPClient *http.Client
		Headers    http.Header
		Token      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    vaultResponse
		wantErr bool
	}{
		{"err", fields{nReq, htcli, nil, ""}, rsp, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				Req:        tt.fields.Req,
				HTTPClient: tt.fields.HTTPClient,
				Headers:    tt.fields.Headers,
				Token:      tt.fields.Token,
			}
			got, err := r.execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
