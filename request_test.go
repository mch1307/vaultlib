package vaultlib

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func Test_newRequest(t *testing.T) {
	testURL, _ := url.Parse("http://localhost:8200")
	emptyReq := new(request)
	type args struct {
		method string
		token  string
		url    *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    *request
		wantErr bool
	}{
		{"err", args{"@", "", testURL}, emptyReq, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRequest(tt.args.method, tt.args.token, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_request_setJSONBody(t *testing.T) {
	var cred AppRoleCredentials
	cred.RoleID = "aa"
	cred.SecretID = "bb"
	htCli := new(http.Client)
	url, _ := url.Parse("http://localhot:8200")
	req, _ := newRequest("GET", "", url)
	ch := make(chan int)

	type fields struct {
		Req        *http.Request
		HTTPClient *http.Client
		Headers    http.Header
		Token      string
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"noJSON", fields{req.Req, htCli, nil, ""}, args{ch}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				Req:        tt.fields.Req,
				HTTPClient: tt.fields.HTTPClient,
				Headers:    tt.fields.Headers,
				Token:      tt.fields.Token,
			}
			if err := r.setJSONBody(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("request.setJSONBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_RawRequest(t *testing.T) {
	_ = os.Unsetenv("VAULT_TOKEN")
	conf := NewConfig()
	conf.AppRoleCredentials.RoleID = vaultRoleID
	conf.AppRoleCredentials.SecretID = vaultSecretID
	vc, err := NewClient(conf)
	if err != nil {
		t.Errorf("Failed to get vault cli %v", err)
	}
	vc.Token.ID = "my-dev-root-vault-token"
	ch := make(chan int)
	type args struct {
		method  string
		path    string
		payload interface{}
	}
	tests := []struct {
		name       string
		cli        *Client
		args       args
		wantResult json.RawMessage
		wantErr    bool
	}{
		{"initEndpoint", vc, args{"GET", "/v1/sys/init", nil}, []byte(`{"initialized":true}
`), false},
		{"badEndpoint", vc, args{"GET", "/v1/wrong/path", nil}, []byte(`{"errors":["no handler for route 'wrong/path'"]}
`), true},
		{"invalidBody", vc, args{"GET", "/v1/sys/init", ch}, nil, true},
		{"noMethod", vc, args{"", "/v1/sys/init", ch}, nil, true},
		{"invalidMethod", vc, args{"@", "/v1/sys/init", ch}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cli
			gotResult, err := c.RawRequest(tt.args.method, tt.args.path, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RawRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotResult), string(tt.wantResult)) {
				t.Errorf("Client.RawRequest() = %v, want %v", string(gotResult), string(tt.wantResult))
			}
		})
	}
}
